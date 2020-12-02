package main
import
(
"log"
"io"
"io/ioutil"
"path/filepath"
"bufio"
"os"
"strings"
"strconv"
"path"
)
func main(){
labelCounterT :=0
labelCounterF :=0
var str string //declaring

reader := bufio.NewReader(os.Stdin) //creating a stdin 
text, err := reader.ReadString('\n') // reading input as string
text = strings.TrimSuffix(text, "\n") //trimming not needed suffix, so file path will be read correctly
text = strings.TrimSuffix(text, "\r") 

str = text
str = filepath.ToSlash(str)
//
dir, err := ioutil.ReadDir(str)
if err != nil{
    log.Fatal(err)
  }
dirName := filepath.Base(str)//getting the current directory name
asmName := dirName + ".asm" // the asm file needs to have the current dir nameerujmi3ttgrg
str= str + "/" //adding slash to the end so it could lead to the file 

file, err := os.Create(str + asmName)// Create- creates a file for reading and writing
if err != nil{
    log.Fatal(err)
  }
wfile, _ := os.OpenFile(str+asmName, os.O_WRONLY|os.O_APPEND, os.ModePerm)// we want to open the file for appending. so if we write new lines the previous ones won't be deleted


for _, f :=range dir { //loop over the directory (the first part _ is not needed because it is the index)
    if!f.IsDir() && filepath.Ext(f.Name()) ==".vm"{ //checking that f is not a directory and then making sure that the type of the file is "VM"
	    file,err = os.Open(str+f.Name())// open the VM file for reading

		reader := bufio.NewReader(file) // getting a file reader
		
		line, _, err := reader.ReadLine()//reads line by line the file 
		s := string(line)//ReadeLine returns []byte but we want the line format to be a string
		if err == io.EOF {//  making sure we are not at the end of the file
				break
		}
		for err == nil{// if this is not null, it means a line can't  be read or that we got to the end of the file.
		  if strings.HasPrefix(s,"push"){//if the begging of the  string is "push"
		    a:=strings.Fields(s)//transforms a string to an array of strings according to space separator
			a = a[1:len(a)] //we want the new array to not contain the first argument
		    if(a[0]== "constant"){//if constant
			   a = a[1:len(a)]
			   wfile.WriteString("@"+a[0])
			   wfile.WriteString("\nD=A\n@SP\nA=M\nM=D\n@SP\nM=M+1\n")
			   }//end constant
			  
			   
			if(a[0]== "pointer"){//if pointer
			   a = a[1:len(a)]
			   int,_ := strconv.Atoi(a[0])//convert from string to int , the function returns two values , the first is number and the second is err 
			   if(int== 0){
			     wfile.WriteString("\n@THIS\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n") 
			     }
			   if(int==1){
			     wfile.WriteString("\n@THAT\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n") 
				 }
			   }//end pointer
			
			 if(a[0]== "that"){//if that
			   a = a[1:len(a)]
			   wfile.WriteString("@"+a[0])
			   wfile.WriteString("\nD=A\n@THAT\nA=M+D\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n")
			   }//end that
			   
             if(a[0]== "this"){//if this
			   a = a[1:len(a)]
			   wfile.WriteString("@"+a[0])
			   wfile.WriteString("\nD=A\n@THIS\nA=M+D\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n")			
		  }//end this
		  
		    if(a[0]== "temp"){//if temp 
			   a = a[1:len(a)]
			   int,_ := strconv.Atoi(a[0])
			   int = int+5
			   wfile.WriteString("@"+strconv.Itoa(int))//Itoa converts int to string, the function WriteSrting  deals only with string type
			   wfile.WriteString("\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n")
			   }//end temp
		  
			 if(a[0]== "local"){//if local 
			   a = a[1:len(a)]
			   wfile.WriteString("@"+a[0])
			   wfile.WriteString("\nD=A\n@LCL\nA=M+D\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n")
			   }//end local
			   
			 if(a[0]== "argument"){//if argument 
			   a = a[1:len(a)]
			   wfile.WriteString("\n@ARG\nD=M\n")
			   wfile.WriteString("@"+ a[0])
			   wfile.WriteString("\nD=D+A\nA=D\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n")
			   }//end argument   
			   
			  if(a[0]== "static"){//if static 
			     a = a[1:len(a)]
				 fileNoExt := strings.TrimSuffix(f.Name(),path.Ext(f.Name()))//we don't want the type of file because according to the hack language it is not needed 
				 wfile.WriteString("\n@"+fileNoExt+"."+a[0])
			     wfile.WriteString("\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n")
			   }//end static 
			   
			   
		 }//end of push
		 
		 
		  if strings.HasPrefix(s,"pop"){
		    a:=strings.Fields(s)
			a = a[1:len(a)]
		    if(a[0]== "pointer"){//if pointer
			  a = a[1:len(a)]
			  int,_ := strconv.Atoi(a[0])
			  if(int== 0){
			     wfile.WriteString("\n@SP\nA=M-1\nD=M\n@THIS\nM=D\n@SP\nM=M-1\n") 
			   }
			   if(int== 1){
			     wfile.WriteString("\n@SP\nA=M-1\nD=M\n@THAT\nM=D\n@SP\nM=M-1\n") 
			   }
			   }//end pointer
			   
			if(a[0]== "this"){//if this
			   a = a[1:len(a)]
			    int,_ := strconv.Atoi(a[0])
			   wfile.WriteString("@SP\nA=M-1\nD=M\n@THIS\nA=M\n")
			   for i:=0; i<int ; i++{
			      wfile.WriteString("A=A+1\n")
			          }
			    wfile.WriteString("M=D\n@SP\nM=M-1\n")
				}//end this	
				
			if(a[0]== "that"){//if that
				a = a[1:len(a)]
				int,_ := strconv.Atoi(a[0])
				wfile.WriteString("@SP\nA=M-1\nD=M\n@THAT\nA=M\n")
			    for i:=0; i<int ; i++{
			      wfile.WriteString("A=A+1\n")
			                        }
			    wfile.WriteString("M=D\n@SP\nM=M-1\n")
		                         }//end that
								 
            if(a[0]== "temp"){//if temp 
			   a = a[1:len(a)]
			   wfile.WriteString("\n@SP\nA=M-1\nD=M\n")
			   int,_ := strconv.Atoi(a[0])
			   int=int+5
			   wfile.WriteString("@"+strconv.Itoa(int))
			   wfile.WriteString("\nM=D\n@SP\nM=M-1\n")
		  }//end temp
	
			if(a[0]== "local"){//if local 
			   a = a[1:len(a)]
			   int,_ := strconv.Atoi(a[0])
				wfile.WriteString("@SP\nA=M-1\nD=M\n@LCL\nA=M\n")
			    for i:=0; i<int ; i++{
			      wfile.WriteString("A=A+1\n")
			                        }
			    wfile.WriteString("M=D\n@SP\nM=M-1\n")
			   }//end local 
			   
			if(a[0]== "argument"){//if argument
			   a = a[1:len(a)]
			    int,_ := strconv.Atoi(a[0])
				wfile.WriteString("@SP\nA=M-1\nD=M\n@ARG\nA=M\n")
			    for i:=0; i<int ; i++{
			      wfile.WriteString("A=A+1\n")
			                        }
			    wfile.WriteString("M=D\n@SP\nM=M-1\n")
			   }//end argument 
			   
			 if(a[0]== "static"){//if static 
			     a = a[1:len(a)]
				 fileNoExt := strings.TrimSuffix(f.Name(),path.Ext(f.Name()))
			     wfile.WriteString("\n@SP\nA=M-1\nD=M\n")
                 wfile.WriteString("\n@"+fileNoExt+"."+a[0])
			     wfile.WriteString("\nM=D\n@SP\nM=M-1\n")
			   }//end static 
																 
			}//end of pop
			
			
		  if strings.HasPrefix(s,"add"){
		    wfile.WriteString("@SP\nA=M-1\nD=M\nA=A-1\nM=D+M\n@SP\nM=M-1\n")
		  }
		  if strings.HasPrefix(s,"sub"){
		    wfile.WriteString("@SP\nA=M-1\nD=M\nA=A-1\nM=M-D\n@SP\nM=M-1\n")
		  }
		  if strings.HasPrefix(s,"neg"){
		    wfile.WriteString("@SP\nA=M-1\nD=M\nM=-D\n")
		  }
		  if strings.HasPrefix(s,"not"){
		    wfile.WriteString("@SP\nA=M-1\nD=M\nM=!D\n")
		  }
		  
		  if strings.HasPrefix(s,"eq"){
		    wfile.WriteString("@SP\nA=M-1\nD=M\nA=A-1\nD=D-M\n@IF_TRUE"+strconv.Itoa(labelCounterT))
			wfile.WriteString("\nD;JEQ\nD=0\n@SP\nA=M-1\nA=A-1\nM=D\n")
			wfile.WriteString("@IF_FALSE"+strconv.Itoa(labelCounterF))
			wfile.WriteString("\n0;JMP\n(IF_TRUE"+strconv.Itoa(labelCounterT))
			wfile.WriteString(")\nD=-1\n@SP\nA=M-1\nA=A-1\nM=D\n")
			wfile.WriteString("(IF_FALSE"+strconv.Itoa(labelCounterF))
			wfile.WriteString(")\n@SP\nM=M-1\n")
			labelCounterT++
			labelCounterF++

		  }
		  if strings.HasPrefix(s,"lt"){
		    wfile.WriteString("@SP\nA=M-1\nD=M\nA=A-1\nD=M-D\n@IF_TRUE"+strconv.Itoa(labelCounterT))
			wfile.WriteString("\nD;JLT\nD=0\n@SP\nA=M-1\nA=A-1\nM=D\n")
			wfile.WriteString("@IF_FALSE"+strconv.Itoa(labelCounterF))
			wfile.WriteString("\n0;JMP\n(IF_TRUE"+strconv.Itoa(labelCounterT))
			wfile.WriteString(")\nD=-1\n@SP\nA=M-1\nA=A-1\nM=D\n")
			wfile.WriteString("(IF_FALSE"+strconv.Itoa(labelCounterF))
			wfile.WriteString(")\n@SP\nM=M-1\n")
			labelCounterT++
			labelCounterF++
			}
		  if strings.HasPrefix(s,"gt"){
		    wfile.WriteString("@SP\nA=M-1\nD=M\nA=A-1\nD=M-D\n@IF_TRUE"+strconv.Itoa(labelCounterT))
			wfile.WriteString("\nD;JGT\nD=0\n@SP\nA=M-1\nA=A-1\nM=D\n")
			wfile.WriteString("@IF_FALSE"+strconv.Itoa(labelCounterF))
			wfile.WriteString("\n0;JMP\n(IF_TRUE"+strconv.Itoa(labelCounterT))
			wfile.WriteString(")\nD=-1\n@SP\nA=M-1\nA=A-1\nM=D\n")
			wfile.WriteString("(IF_FALSE"+strconv.Itoa(labelCounterF))
			wfile.WriteString(")\n@SP\nM=M-1\n")
			labelCounterT++
			labelCounterF++
			}
		  if strings.HasPrefix(s,"and"){
		    wfile.WriteString("@SP\nA=M-1\nD=M\nA=A-1\nM=M&D\n@SP\nM=M-1\n")
			}
		  if strings.HasPrefix(s,"or"){
		    wfile.WriteString("@SP\nA=M-1\nD=M\nA=A-1\nM=M|D\n@SP\nM=M-1\n")
		  }
		  
		  line, _,err = reader.ReadLine()// reader is stdin reader created with bufio library, ReadLine returns 3 values: an array of bytes, boolean that set if the buffer was small for the amount of data, and err type
		  s = string(line) // converting the array of bytes to string
		}
		
    }
	
 }  //end of the outer format
 wfile.Close()
 }

