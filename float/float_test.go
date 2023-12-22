package float

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/test"
)

type F32Circuit struct {
	X  frontend.Variable `gnark:",secret"`
	Y  frontend.Variable `gnark:",secret"`
	Z  frontend.Variable `gnark:",public"`
	op string
}

func (c *F32Circuit) Define(api frontend.API) error {
	ctx := NewContext(api, 8, 23)
	x := ctx.NewFloat(c.X)
	y := ctx.NewFloat(c.Y)
	z := ctx.NewFloat(c.Z)
	ctx.AssertIsEqual(reflect.ValueOf(&ctx).MethodByName(c.op).Call([]reflect.Value{reflect.ValueOf(x), reflect.ValueOf(y)})[0].Interface().(FloatVar), z)
	return nil
}

type F64Circuit struct {
	X  frontend.Variable `gnark:",secret"`
	Y  frontend.Variable `gnark:",secret"`
	Z  frontend.Variable `gnark:",public"`
	op string
}

func (c *F64Circuit) Define(api frontend.API) error {
	ctx := NewContext(api, 11, 52)
	x := ctx.NewFloat(c.X)
	y := ctx.NewFloat(c.Y)
	z := ctx.NewFloat(c.Z)
	ctx.AssertIsEqual(reflect.ValueOf(&ctx).MethodByName(c.op).Call([]reflect.Value{reflect.ValueOf(x), reflect.ValueOf(y)})[0].Interface().(FloatVar), z)
	return nil
}

func TestF32Circuit(t *testing.T) {
	assert := test.NewAssert(t)

	ops := []string{"Add", "Sub", "Mul", "Div"}

	for _, op := range ops {
		path, _ := filepath.Abs(fmt.Sprintf("../data/f32/%s", strings.ToLower(op)))
		file, _ := os.Open(path)
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			data := strings.Fields(scanner.Text())
			a, _ := new(big.Int).SetString(data[0], 16)
			b, _ := new(big.Int).SetString(data[1], 16)
			c, _ := new(big.Int).SetString(data[2], 16)

			assert.ProverSucceeded(
				&F32Circuit{X: 0, Y: 0, Z: 0, op: op},
				&F32Circuit{X: a, Y: b, Z: c, op: op},
				test.WithCurves(ecc.BN254),
				test.WithBackends(backend.GROTH16, backend.PLONK),
			)
		}
	}
}

func TestF64Circuit(t *testing.T) {
	assert := test.NewAssert(t)

	ops := []string{"Add", "Sub", "Mul", "Div"}

	for _, op := range ops {
		path, _ := filepath.Abs(fmt.Sprintf("../data/f64/%s", strings.ToLower(op)))
		file, _ := os.Open(path)
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			data := strings.Fields(scanner.Text())
			a, _ := new(big.Int).SetString(data[0], 16)
			b, _ := new(big.Int).SetString(data[1], 16)
			c, _ := new(big.Int).SetString(data[2], 16)

			assert.ProverSucceeded(
				&F64Circuit{X: 0, Y: 0, Z: 0, op: op},
				&F64Circuit{X: a, Y: b, Z: c, op: op},
				test.WithCurves(ecc.BN254),
				test.WithBackends(backend.GROTH16, backend.PLONK),
			)
		}
	}
}
