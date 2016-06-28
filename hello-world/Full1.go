package main

import (
	"./lexport"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Device struct {
	Name, Vendor string
}

// INTERFACES!!!
type IDeviceRecords interface {
	AddDevice(d *Device)
	GetDevices()
}

type DeviceRecords struct {
	Arrd []*Device
}

func (dv *DeviceRecords) AddDevice(d *Device) {
	dv.Arrd = append(dv.Arrd, d)
}

func (dv *DeviceRecords) GetDevices() {
	fmt.Println(len(dv.Arrd))
	for i, v := range dv.Arrd {
		fmt.Println(i, v.Vendor)
	}
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered in controller, panic was: %q\n", r)
		}
	}()

	// using printf
	var x string = "yellow"
	fmt.Printf("Something %d, %s\n", lexport.Exportable(), x)

	// lexport.unexportable()

	// arrays
	var arr [4]string
	arr[0] = "a"
	fmt.Println(arr[0], "\n")

	// slices
	slice1 := make([]string, 0, 0)
	slice1 = append(slice1, "A", "B")
	fmt.Println(slice1)
	slice2 := slice1[0:1]
	fmt.Println(slice2)

	//maps
	map1 := make(map[int]string)
	map1[5] = "prg"

	// and channels
	for i := 0; i < 10; i++ {
		map1[rand.Int()%50] = "xyz"
	}

	// fmt.Println(map1)
	for k, v := range map1 {
		fmt.Println("Map: ", k, v)
	}

	// infinite loops
	for /*true*/ {
		break
	}

	// while emulation
	var u int = 5
	for u < 10 {
		fmt.Println(u)
		u++
	}

	// ranges
	var slc []string = []string{"a", "b", "c"} // decl syntax
	// slc := []string{"a", "b", "c"} // equiv decl syntax

	for i, v := range slc {
		fmt.Printf("%d, %s\n", i, v) // can use semicolon
	}

	fmt.Println(len(map1), "<=>", 10)
	if len(map1) < 9 {
		fmt.Println("X")
	} else {
		fmt.Println("Y")
	}

	// structs and things
	type Something struct {
		Nodal    bool
		Distance int64
	}

	a := Something{
		Nodal:    true,
		Distance: 5,
	}
	fmt.Println(a)

	// modify byref, ptrs
	ivyx := 5
	modByRef(&ivyx)
	fmt.Println(ivyx)

	//new()?!
	var intPr *int = new(int)
	fmt.Println(intPr) // addr

	// find keys in map
	fmt.Println("\n\n\n\n\n")
	if v, ok := map1[5]; ok {
		fmt.Println(v)
	} else {
		fmt.Println("Lookup Failed")
	}

	// variadic
	fmt.Println(variadics(5, 5, 7, 7))

	// switch/case
	var d Device = Device{
		Name:   "quux",
		Vendor: "foox",
	}
	fmt.Println(checkDevice(&d))

	// errors
	if i := genError(3); i != nil {
		fmt.Println(i)
	}

	if res, err := Divide(5, 3); err == nil {
		fmt.Println(res)
	} else {
		fmt.Println(err)
	}

	// first class functions.
	var fcf func() int = func() int {
		return 42
	}
	fmt.Println("Ans: ", fcf())

	// decorator pattern
	decorator()

	// interfaces
	var dr IDeviceRecords
	dr = &(DeviceRecords{})
	dr.AddDevice(&d)
	dr.GetDevices()

	// concurrency
	for i := 0; i < 5; i++ {
		go func() {
			fmt.Println("Cheap threads!")
			time.Sleep(100 * time.Millisecond)
		}()
	}
	select {}

	// channels (stream-like)
	// ch := make(chan bool, 2)
	// fmt.Println(<-ch)
	// ch <- true

	// panic
	panic("This is panic by design")
	fmt.Println("post-panic")

}

func decorator() {
	defer fmt.Println("post-decorating")
	fmt.Println("DecoratorBody")
}

func genError(i int) error {
	if i < 5 {
		return errors.New("FAIL")
	}
	return nil
}

func modByRef(a *int) {
	*a = 10
}

func variadics(inp ...int) int {
	total := 0
	for _, v := range inp {
		total += v
	}
	return total
}

func Divide(num int, denom int) (float64, error) {
	if denom == 0 {
		return 0.0, errors.New("DIV-BY-ZERO")
	} else {
		return float64(num) / float64(denom), nil
	}
}

func checkDevice(d *Device) bool {
	switch nm := d.Vendor; nm {
	case "blah":
		return true
	case "foo":
		return true
	default:
		return false
	}
}
