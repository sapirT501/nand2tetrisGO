package main
import (
"bufio"
"os"
"targil5/Tokenizing"
"targil5/targil5"
)


func main(){
    reader := bufio.NewReader(os.Stdin) //creating a stdin 
   text, _ := reader.ReadString('\n')
   Tokenizing.Main(text)
   targil5.Main(text)
}