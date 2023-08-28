package interpreter

import "io"

func asBin(operations *[]Operation, writer io.ByteWriter) error {
    writer.WriteByte(0xCA)
    writer.WriteByte(0xFE)
    writer.WriteByte(0xBA)
    writer.WriteByte(0xBE)




    return nil
}
