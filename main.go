package main

import (
    "encoding/binary"
    "log"
    "modbus-demo/modbus"
    "os"
    "time"
)

func main() {

    handler := modbus.NewTCPClient2Handler("localhost:502")
    handler.Timeout = 10 * time.Second
    handler.SlaveId = 0xFF
    handler.Logger = log.New(os.Stdout, "test: ", log.LstdFlags)
    err := handler.Connect()
    if err != nil {
        panic(err)
    }

    handler.SlaveId = 1
    defer handler.Close()
    client := modbus.NewClient(handler)
    //coils, err := client.ReadCoils(50, 1)

    //client.MaskWriteRegister(40001, 1,182 )
    client.WriteSingleRegister(40001, 200 )
    //time.Sleep(time.Millisecond * 500)
    registers, err := client.ReadHoldingRegisters(40001, 1)
    if err != nil {
        panic(err)
    }
    println(binary.BigEndian.Uint16(registers))
    time.Sleep(time.Millisecond * 1000)
    registers, err = client.ReadHoldingRegisters(40001, 1)
    if err != nil {
        panic(err)
    }
    println(binary.BigEndian.Uint16(registers))


    //var num = 3
    //registers, err := client.ReadHoldingRegisters(4, uint16(num))
    //
    //if err != nil {
    //    panic(err)
    //}
    //for i := 0; i < num; i++ {
    //    println(binary.BigEndian.Uint16(registers[i * 2: i * 2 + 2]))
    //}
    //
    //num2 := 50
    //for i:=0; i< 1; i ++ {
    //    registers, err := client.ReadHoldingRegisters(1, uint16(num2))
    //    if err != nil {
    //        panic(err)
    //    }
    //    for i := 0; i < num2; i++ {
    //        print(binary.BigEndian.Uint32(registers[i * 4: i * 4 + 4]), ",")
    //    }
    //    println()
    //}




}
