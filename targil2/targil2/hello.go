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
funcCounter :=0

var str string //declaring

reader := bufio.NewReader(os.Stdin) //creating a stdin 
text, err := reader.ReadString('\n') // reading input as string
text = strings.TrimSuffix(text, "\n") //triming not needed suffix, so filepath will be read correctly
text = strings.TrimSuffix(text, "\r") 

str = text
str = filepath.ToSlash(str)
//
dir, err := ioutil.ReadDir(str)
if err != nil{
    log.Fatal(err)
  }
dirName := filepath.Base(str)//getting the current dir name
asmName := dirName + ".asm" // the asm file needs to have the current dir name
str= str + "/" //adding slash to the end so it could lead to the file 

file, err := os.Create(str + asmName)// Create- creates a file for read and write
if err != nil{
    log.Fatal(err)
  }
wfile, _ := os.OpenFile(str+asmName, os.O_WRONLY|os.O_APPEND, os.ModePerm)// we want to openf for appending so if we write new lines the previous wont be deleted
for _,num :=range dir{
 if(num.Name() == "Sys.vm"){  
    wfile.WriteString("@256\nD=A\n@SP\nM=D\n")
    wfile.WriteString("@Sys.init.returnAdd\nD=A\n@SP\nA=M\nM=D\n@SP\nM=M+1\n@LCL\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n")
    wfile.WriteString("@ARG\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n@THIS\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n@THAT\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n")
    wfile.WriteString("@SP\nD=M\n@5\nD=D-A\n@ARG\nM=D\n@SP\nD=M\n@LCL\nM=D\n@Sys.init\n0;JMP\n(Sys.init.returnAdd)\n")
    }
}
for _, f :=range dir { //loop over dir (the first part _ not neaded beacause it is the index)
    if!f.IsDir() && filepath.Ext(f.Name()) ==".vm"{ //checking that f is mot a directory and than making sure that the type of the file is "vm"
	    file,err = os.Open(str+f.Name())// open the vm file for reading
		reader := bufio.NewReader(file) // getting a file reader
		line, _, err := reader.ReadLine()//reads line by line the file 
		s := string(line)// readeline returns []byte but we want the line format to be string
		fileNoExt := strings.TrimSuffix(f.Name(),path.Ext(f.Name()))//we dont want the type of file because according to the hack langauage it is not needed 
		if err == io.EOF {// making sure we are not in the end of file
				break
		}
		for err == nil{// if this is not null it means a line cant be read or we got to the end 
		  if strings.HasPrefix(s,"push"){//if the begging of the  string is "push"
		    a:=strings.Fields(s)//transforms a string to array of strings according to space separator
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
			   wfile.WriteString("@"+strconv.Itoa(int))//Itoa converts int to string, the function writesrting  deals only with string type
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
			
			
		  if strings.HasPrefix(s,"function"){
		    a:=strings.Fields(s)
			a = a[1:len(a)]
			wfile.WriteString("("+ a[0] +")\n")
			a = a[1:len(a)]
			int,_ := strconv.Atoi(a[0])
			
			
			for i:=0; i<int ; i++{
			      wfile.WriteString("\n@0\nD=A\nA=M\nM=D\n@SP\nM=M+1\n") 
			   }
               }//end function

          if strings.HasPrefix(s,"return"){ //if return
		     wfile.WriteString("\n@LCL\nD=M\n") // FRAME = LCL
			 wfile.WriteString("@5\nA=D-A\nD=M\n@13\nM=D\n") // RET = * (FRAME-5) // RAM[13] = (LOCAL - 5)
             wfile.WriteString("@SP\nM=M-1\nA=M\nD=M\n@ARG\nA=M\nM=D\n")// * ARG = pop()
			 wfile.WriteString("@ARG\nD=M\n@SP\nM=D+1\n")// SP = ARG+1 
			 wfile.WriteString("@LCL\nM=M-1\nA=M\nD=M\n@THAT\nM=D\n")// THAT = *(FRAM-1)
             wfile.WriteString("@LCL\nM=M-1\nA=M\nD=M\n@THIS\nM=D\n")// THIS = *(FRAM-2)
             wfile.WriteString("@LCL\nM=M-1\nA=M\nD=M\n@ARG\nM=D\n")// ARG = *(FRAM-3)
			 wfile.WriteString("@LCL\nM=M-1\nA=M\nD=M\n@LCL\nM=D\n")//  LCL= *(FRAM-4)
			 wfile.WriteString("@13\nA=M\n0;JMP\n")// goto RET
               }//end return			   
			
		  if strings.HasPrefix(s,"call"){//if call
		     a:=strings.Fields(s)
		     a = a[1:len(a)]
			 stringFCounter:= strconv.Itoa(funcCounter)
             wfile.WriteString("@"+a[0]+".ReturnAddress"+stringFCounter+"\n") // push return-address
  			 wfile.WriteString("D=A\n@SP\nA=M\nM=D\n@SP\nM=M+1\n")	 
			 wfile.WriteString("@LCL\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n")//push LCL
             wfile.WriteString("\n@ARG\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n")//push ARG
			 wfile.WriteString("@THIS\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n")//push THIS
			 wfile.WriteString("@THAT\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n")//push THAT
			 FuncName := a[0]
			 a = a[1:len(a)]
			 int,_ := strconv.Atoi(a[0])
			 int = int + 5
			 wfile.WriteString("@SP\nD=M\n@"+strconv.Itoa(int)+"\nD=D-A\n@ARG\nM=D\n")// ARG = SP-n-5 
			 wfile.WriteString("@SP\nD=M\n@LCL\nM=D\n")// LCL = SP
			 wfile.WriteString("@"+FuncName+"\n0;JMP\n")// goto FuncName
                
			 wfile.WriteString("("+FuncName+".ReturnAddress"+stringFCounter+")\n")// label return-address
			 funcCounter++
			 
		  }		//end call
		  
	      if strings.HasPrefix(s,"label"){//if label
		     a:=strings.Fields(s)
		     a = a[1:len(a)]
		     fileNoExt := strings.TrimSuffix(f.Name(),path.Ext(f.Name()))
			 wfile.WriteString("("+fileNoExt+"."+a[0]+")\n")
		  }		//end label
		  
		   if strings.HasPrefix(s,"goto"){//if goto
		     a:=strings.Fields(s)
			 a = a[1:len(a)]
		     fileNoExt := strings.TrimSuffix(f.Name(),path.Ext(f.Name()))
			 wfile.WriteString("@"+fileNoExt+"."+a[0])
			 wfile.WriteString("\n0;JMP\n")
		  }		//end goto
		  
		  if strings.HasPrefix(s,"if-goto"){//if-goto
		     a:=strings.Fields(s)
			 a = a[1:len(a)]
			 wfile.WriteString("\n@SP\nA=M-1\nD=M\n@SP\nM=M-1\n")//load the top stack to D and immidiately decressing SP beacause after JNE it is not possible
			 wfile.WriteString("@"+fileNoExt+"."+a[0])
			 wfile.WriteString("\nD;JNE\n")  
		  }		//end if-goto
    		  
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
		  
		  line, _,err = reader.ReadLine()// reader is stdin reader created with bufio labarary, ReadLine returns 3 values : array of byts , boolian that set if the buffer was small for the amount of data, and err type
		  s = string(line) // converting the array of bytes to string
		}
		
    }
	
 } // end of outer format
 wfile.Close()
 }

