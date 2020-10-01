package main
import (
"bufio"
"os"
"targil42/Tokenizing"
"targil42/Parsing"
)


func main(){
    reader := bufio.NewReader(os.Stdin) //creating a stdin 
   text, _ := reader.ReadString('\n')
   Tokenizing.Main(text)
   Parsing.Main(text)
}