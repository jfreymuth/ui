package gldraw

import (
	"fmt"
	"image"
	"os"
	"strings"

	"github.com/jfreymuth/ui/draw"

	"github.com/go-gl/gl/v3.3-core/gl"
	m "github.com/go-gl/mathgl/mgl32"
)

type buffer struct {
	vao, vbo uint32
	buf      []vertex
	program  uint32
	sizeLoc  int32
}

func (b *buffer) init(cap int) {
	b.buf = make([]vertex, 0, cap)
	b.program = createProgram(vss, fss)
	gl.UseProgram(b.program)
	b.sizeLoc = gl.GetUniformLocation(b.program, gl.Str("screenSize\x00"))
	gl.GenBuffers(1, &b.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, b.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, vertexSize*cap, nil, gl.STREAM_DRAW)
	gl.GenVertexArrays(1, &b.vao)
	gl.BindVertexArray(b.vao)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.INT, false, vertexSize, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, vertexSize, gl.PtrOffset(8))
	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointer(2, 4, gl.UNSIGNED_BYTE, true, vertexSize, gl.PtrOffset(16))
}

func (b *buffer) setScreenSize(w, h int) {
	gl.UseProgram(b.program)
	gl.Uniform2f(b.sizeLoc, float32(w), float32(h))
}

func (b *buffer) allocate(n int) []vertex {
	if n > b.free() {
		b.flush()
	}
	i := len(b.buf)
	b.buf = b.buf[:i+n]
	return b.buf[i : i+n]
}

func (b *buffer) free() int {
	return cap(b.buf) - len(b.buf)
}

func (b *buffer) rect(r image.Rectangle, tmin, tmax m.Vec2, color draw.Color) {
	if r.Empty() {
		return
	}
	v := b.allocate(6)
	v[0] = vertex{[2]int32{int32(r.Min.X), int32(r.Min.Y)}, tmin, color}
	v[1] = vertex{[2]int32{int32(r.Min.X), int32(r.Max.Y)}, m.Vec2{tmin[0], tmax[1]}, color}
	v[2] = vertex{[2]int32{int32(r.Max.X), int32(r.Min.Y)}, m.Vec2{tmax[0], tmin[1]}, color}
	v[3] = v[2]
	v[4] = vertex{[2]int32{int32(r.Max.X), int32(r.Max.Y)}, tmax, color}
	v[5] = v[1]
}

func (b *buffer) flush() {
	if len(b.buf) == 0 {
		return
	}
	gl.UseProgram(b.program)
	gl.BindBuffer(gl.ARRAY_BUFFER, b.vbo)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(b.buf)*vertexSize, gl.Ptr(b.buf))
	gl.BindVertexArray(b.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(b.buf)))
	b.buf = b.buf[:0]
}

type vertex struct {
	pos   [2]int32
	tc    m.Vec2
	color [4]uint8
}

const vertexSize = 20

func createProgram(vss, fss string) uint32 {
	program := gl.CreateProgram()
	shaders := []struct {
		source string
		typ    uint32
	}{{vss, gl.VERTEX_SHADER}, {fss, gl.FRAGMENT_SHADER}}
	for _, shader := range shaders {
		s := gl.CreateShader(shader.typ)
		csource, length := gl.Str(shader.source), int32(len(shader.source))
		gl.ShaderSource(s, 1, &csource, &length)
		gl.CompileShader(s)

		var logLen int32
		gl.GetShaderiv(s, gl.INFO_LOG_LENGTH, &logLen)
		if logLen > 0 {
			log := strings.Repeat("\x00", int(logLen))
			gl.GetShaderInfoLog(s, logLen, nil, gl.Str(log))
			log = log[:len(log)-1]
			fmt.Printf("GL shader compilation:\n%s\n", log)
		}

		var status int32
		gl.GetShaderiv(s, gl.COMPILE_STATUS, &status)
		if status == gl.FALSE {
			gl.DeleteShader(s)
			fmt.Println("GL shader compilation failed")
			os.Exit(1)
		}
		gl.AttachShader(program, s)
	}
	gl.LinkProgram(program)

	var logLength int32
	gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)
	if logLength > 0 {
		log := strings.Repeat("\x00", int(logLength))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))
		log = log[:len(log)-1]
		fmt.Printf("GL shader compilation:\n%s\n", log)
	}

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		gl.DeleteProgram(program)
		fmt.Println("GL shader compilation failed")
		os.Exit(1)
	}

	return program
}

var vss = `#version 330
layout(location = 0) in vec2 pos;
layout(location = 1) in vec2 tc;
layout(location = 2) in vec4 color;

out vec2 ftc;
out vec4 fcol;

uniform vec2 screenSize;

void main() {
	gl_Position = vec4(vec2(-1, 1) + pos/screenSize*vec2(2, -2), 0, 1);
	ftc = tc;
	fcol = color;
}
` + "\x00"

var fss = `#version 330
in vec2 ftc;
in vec4 fcol;

out vec4 outcol;

uniform sampler2D image;

void main() {
	outcol = fcol * texture(image, ftc);
}
` + "\x00"
