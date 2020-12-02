package Tokenizing

import (
	"bufio"
	"io"
	"log"
	"strings"
"io/ioutil"
"path/filepath"
"os"
"path"
//"strconv"
)


func keyWord(word string) bool{
	if(word=="class" || word=="constructor" || word=="function" || word=="method" || word=="field" || word=="static" || word=="var" || word=="int" || word=="char" || word=="boolean" || word=="void" || word=="true" || word=="false" || word=="null" || word=="this" || word=="let" || word=="do" || word=="if" || word=="else" || word=="while" || word=="return"){
			return true
		}
	return false
}

func Main(text string) {
    var str string //declaring

   //reader := bufio.NewReader(os.Stdin) //creating a stdin 
   //text, _ := reader.ReadString('\n') // reading input as string
   text = strings.TrimSuffix(text, "\n") //trimming not needed suffix, so filepath will be read correctly
   text = strings.TrimSuffix(text, "\r") 

   str = text
   str = filepath.ToSlash(str)

//
  dir, err := ioutil.ReadDir(str)
if err != nil{
    log.Fatal(err)
  }

 
str= str + "/" //adding a slash to the end so it could lead to the file 
t:=""


for _, f :=range dir { //loop over dir (the first part _ not needed because it is the index)
    if!f.IsDir() && filepath.Ext(f.Name()) ==".jack"{ //checking that f is mot a directory and then making sure that the type of the file is "VM"
		file,_ := os.Open(str+f.Name())// open the vm file for reading
		r := bufio.NewReader(file) // getting a file reader
		fileNoExt := strings.TrimSuffix(f.Name(),path.Ext(f.Name()))//we don't want the type of file because according to the hack language it is not needed 
		
		/**creating XML file**/
		
		xmlName := fileNoExt +"T"+ ".xml"
		file, err = os.Create(str + xmlName)// Create- creates a file for reading and writing
        if err != nil{
           log.Fatal(err)
           }
		wfile, _ := os.OpenFile(str+xmlName, os.O_WRONLY|os.O_APPEND, os.ModePerm)// we want to open the file for appending so if we write new lines the previous won't be deleted   
		wfile.WriteString("<tokens>\n")
	for {//q0
	    isCorrect:=true
		word:=""
		elseif:=true
		if c, _, err := r.ReadRune(); err != nil {   //ReadRune read the file char by char
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		} else {
		    t=string(c)
			if !(t=="\n"||t==" ") && !strings.ContainsAny(t,"|") && (t=="_" || strings.ContainsAny(t,"a | b | c | d  | e | f | g | h | i | j | k | l | m | n | o | p | q | r | s | t | u | v | w | x | y | z | A | B | C | D  | E | F | G | H | I | J | K | L | M | N | O | P | Q | R | S | T | U | V | W | X | Y | Z" )){//identifier
					    
						word=word+t
						for{//q1 loop
						c, _, err := r.ReadRune()
						t=string(c)
                        if  err != nil {
						   
			              if err == io.EOF {
				           break
			              } else {
				           log.Fatal(err)
			               }
						   
					    }else{
						   if (t=="\n"||t==" "){
						   isCorrect=false
						   break}
						   if (t=="_" || strings.ContainsAny(t,"0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9") || strings.ContainsAny(t,"a | b | c | d  | e | f | g | h | i | j | k | l | m | n | o | p | q | r | s | t | u | v | w | x | y | z | A | B | C | D  | E | F | G | H | I | J | K | L | M | N | O | P | Q | R | S | T | U | V | W | X | Y | Z" )){
							  word=word+t
						    }else{break}
						}
						}
						if(keyWord(word)){
						   
						   wfile.WriteString("<keyword>"+word+ "</keyword>\n")
						   
						  }
						if(!keyWord(word)){
						    
						   wfile.WriteString("<identifier>" + word + "</identifier>\n")
							
						}
					
					}//end q1
						
					
					
					if !(t=="\n"||t==" ")&&!strings.ContainsAny(t,"|") &&( isCorrect && strings.ContainsAny(t,"0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9")){//integerConstant
					   	word=word+t
					     for{//q2 loop
						 c, _, err := r.ReadRune()
						 t=string(c)
                        if  err != nil {
			              if err == io.EOF {
				           break
			              }else{
				           log.Fatal(err)
			               }
					       }else{
						    if (strings.ContainsAny(t,"0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9") ){
							  word=word+t
						    }else{ 
							break}
						  }
						} //qNumber  
						
						if(isCorrect){
						    //integerConstantonst
							wfile.WriteString("<integerConstant>"+ word +"</integerConstant>\n")
						}
			
			        }
					
					
				/********this function check if is comment and skipping this*********/	
			   if(c=='/'){
					c, _, _ =r.ReadRune()
					if(c=='/'){
					isCorrect=false
					r.ReadLine()
					}
					if(c=='*'){
						for{// loop to find *
							if c, _, err = r.ReadRune(); err != nil {
								if err == io.EOF {
									break
							} else {
								log.Fatal(err)
								}
							}else{
						    
						   if (c=='*'){
							   c, _, _ =r.ReadRune()
							   if(c=='/'){isCorrect=false
							   break}
						    }
						}
						}
				
				}
				}
				
				/********check if this symbol*********/
					if( !(t=="\n"||t==" ") && isCorrect && ( strings.ContainsAny(t,"|") || strings.ContainsAny(t,"{ | } | ( | )  | [ | ] | . | , | ; | + | - | * | / | & |  < | > | = | ~ "))){//symbol
						     
						          if(elseif && t=="<"){
									wfile.WriteString("<symbol>&lt;</symbol>\n")
									elseif=false
									}
									if(t==">"){
									wfile.WriteString("<symbol>&gt;</symbol>\n")
									elseif=false
									}
									if(t=="&"){
									wfile.WriteString("<symbol>&amp;</symbol>\n")
									elseif=false
									}
									if(elseif){
								     wfile.WriteString("<symbol>"+ t +"</symbol>\n")}
						   
						}
						
						
				    /********check if this begins with the quote*********/
					if(  c=='"' ){//stringConstant
						
						for{// loop
                        if c, _, err := r.ReadRune(); err != nil {
			              if err == io.EOF {
				           break
			              } else {
				           log.Fatal(err)
			               }
					    }else{
						    t=string(c)
						   if (c!='"'){
							  word=word+t
						    }else{ 
							break}
						}
						}
						   
						   if(isCorrect){
								wfile.WriteString("<stringConstant>"+word+"</stringConstant>\n")
								}
						}//qStringCon
			 }
		    }
			wfile.WriteString("</tokens>")
			wfile.Close()
		}
			 	
	}
}

