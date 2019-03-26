package gldraw

import (
	"image"

	"github.com/jfreymuth/ui/draw"

	"github.com/go-gl/gl/v3.3-core/gl"
	m "github.com/go-gl/mathgl/mgl32"
)

type sbuffer struct {
	vao, vbo uint32
	buf      []svertex
	program  uint32
	sizeLoc  int32
}

func (b *sbuffer) init(cap int) {
	b.buf = make([]svertex, 0, cap)
	b.program = createProgram(svss, sfss)
	gl.UseProgram(b.program)
	b.sizeLoc = gl.GetUniformLocation(b.program, gl.Str("screenSize\x00"))
	gl.GenBuffers(1, &b.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, b.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, svertexSize*cap, nil, gl.STREAM_DRAW)
	gl.GenVertexArrays(1, &b.vao)
	gl.BindVertexArray(b.vao)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, svertexSize, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointer(2, 4, gl.UNSIGNED_BYTE, true, svertexSize, gl.PtrOffset(12))
	gl.EnableVertexAttribArray(3)
	gl.VertexAttribPointer(3, 4, gl.INT, false, svertexSize, gl.PtrOffset(16))
}

func (b *sbuffer) setScreenSize(w, h int) {
	gl.UseProgram(b.program)
	gl.Uniform2f(b.sizeLoc, float32(w), float32(h))
}

func (b *sbuffer) allocate(n int) []svertex {
	if n > b.free() {
		b.flush()
	}
	i := len(b.buf)
	b.buf = b.buf[:i+n]
	return b.buf[i : i+n]
}

func (b *sbuffer) free() int {
	return cap(b.buf) - len(b.buf)
}

func (b *sbuffer) rect(min, max m.Vec2, rect image.Rectangle, r float32, color draw.Color) {
	if rect.Empty() {
		return
	}
	rect32 := [4]int32{int32(rect.Min.X), int32(rect.Min.Y), int32(rect.Max.X), int32(rect.Max.Y)}
	v := b.allocate(6)
	v[0] = svertex{m.Vec3{min[0] - r, min[1] - r, r * .4}, color, rect32}
	v[1] = svertex{m.Vec3{min[0] - r, max[1] + r, r * .4}, color, rect32}
	v[2] = svertex{m.Vec3{max[0] + r, min[1] - r, r * .4}, color, rect32}
	v[3] = v[2]
	v[4] = svertex{m.Vec3{max[0] + r, max[1] + r, r * .4}, color, rect32}
	v[5] = v[1]
}

func (b *sbuffer) flush() {
	if len(b.buf) == 0 {
		return
	}
	gl.UseProgram(b.program)
	gl.BindBuffer(gl.ARRAY_BUFFER, b.vbo)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(b.buf)*svertexSize, gl.Ptr(b.buf))
	gl.BindVertexArray(b.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(b.buf)))
	b.buf = b.buf[:0]
}

type svertex struct {
	pos   m.Vec3
	color [4]uint8
	rect  [4]int32
}

const svertexSize = 32

var svss = `#version 330
layout(location = 0) in vec3 pos;
layout(location = 2) in vec4 color;
layout(location = 3) in vec4 rect;

out vec4 fcol;
out vec2 fpos;
out vec4 frect;
out float fr;

uniform vec2 screenSize;

void main() {
	gl_Position = vec4(vec2(-1, 1) + pos.xy/screenSize*vec2(2, -2), 0, 1);
	fpos = pos.xy;
	fr = pos.z;
	frect = rect;
	fcol = color;
}
` + "\x00"

var sfss = `#version 330
in vec4 fcol;
in vec2 fpos;
in vec4 frect;
in float fr;

out vec4 outcol;

uniform sampler2D image;

vec4 erf(vec4 x) {
	vec4 s = sign(x), a = abs(x);
	x = 1.0 + (0.278393 + (0.230389 + 0.078108 * (a * a)) * a) * a;
	x *= x;
	return s - s / (x * x);
}

float boxShadow(vec2 lower, vec2 upper, vec2 point, float sigma) {
	vec4 query = vec4(point - lower, point - upper);
	vec4 integral = 0.5 + 0.5 * erf(query * (sqrt(0.5) / sigma));
	return (integral.z - integral.x) * (integral.w - integral.y);
}

void main() {
	outcol = fcol * boxShadow(frect.xy, frect.zw, fpos, fr);
}
` + "\x00"
