package Parsing
import
(
"log"
"io"
"io/ioutil"
"path/filepath"
"bufio"
"os"
"strings"
//"strconv"

)


func Main(text string){
//reader := bufio.NewReader(os.Stdin) //creating a stdin 
//text, err := reader.ReadString('\n') // reading input as string
text = strings.TrimSuffix(text, "\n") //triming not needed suffix, so filepath will be read correctly
text = strings.TrimSuffix(text, "\r") 

str := text
str = filepath.ToSlash(str)

//
dir, err := ioutil.ReadDir(str)
if err != nil{
    log.Fatal(err)
  }
//dirName := filepath.Base(str)//getting the current dir name
 
str= str + "/" //adding slash to the end so it could lead to the file 



for _, f :=range dir { //loop over dir (the first part _ not neaded beacause it is the index)
    fileName:=f.Name()
	if!f.IsDir() && strings.HasSuffix(fileName,"T.xml"){ //checking that f is mot a directory and than making sure that the type of the file is "vm"
	    
		file,err := os.Open(str+f.Name())// open the TOKEN file for reading
		reader := bufio.NewReaderSize(file,20000) // getting a file reader
				
		fileNoExt := strings.TrimSuffix(f.Name(),"T.xml")//we dont want the type of file because according to the hack langauage it is not needed 
		
		/**creating xml file**/
		
		xmlName := fileNoExt +".xml"
		file, err = os.Create(str + xmlName)// Create- creates a file for read and write
        if err != nil{
           log.Fatal(err)
           }
		wfile, _ := os.OpenFile(str+xmlName, os.O_WRONLY|os.O_APPEND, os.ModePerm)// we want to openf for appending so if we write new lines the previous wont be deleted

		/***/
		
		if err == io.EOF {// making sure we are not in the end of file
				break
		}
        ParseClass(reader ,wfile)
        




}//end of if

}//end of for dir


}//end of main



/****getNextToken****/
func getNextToken(reader *bufio.Reader) (string,*bufio.Reader){
	line, _, _ := reader.ReadLine()//reads line by line the file 
	s := string(line)// readeline returns []byte but we want the line format to be string
	return s,reader
	
}

/***func CheckNextToken***/
func CheckNextToken (reader bufio.Reader,s1 string) bool{
	line, _, _ := reader.ReadLine()//reads line by line the file line is [] byte
	s:= string(line)
	s= strings.TrimSuffix(s, "\n") //triming not needed suffix, so filepath will be read correctly
	s= strings.TrimSuffix(s, "\r")
	if(strings.Contains(s,s1)){
		return true	
	}
	return false
}

/***func CheckNextToken2***/
func CheckNextToken2 (reader bufio.Reader,s1 string) bool{
	line, _, _ := reader.ReadLine()//reads line by line the file line is [] byte
	s:= string(line)
	line, _, _ = reader.ReadLine()//reads line by line the file line is [] byte
	s= string(line)
	
	
	if(strings.ContainsAny(s,s1)){
		return true	
	}
	return false
}



/***********Class***********/

func ParseClass(reader *bufio.Reader ,wfile *os.File){
    reader.ReadLine() 	
	wfile.WriteString("<class>\n")
    s:=""
	
	s,reader=getNextToken(reader) 
	wfile.WriteString(s+"\n")//class
	
	s,reader=getNextToken(reader) 
	wfile.WriteString(s+"\n")//type
	
	s,reader=getNextToken(reader)
	wfile.WriteString(s+"\n")//{
	
	for(CheckNextToken((*reader),"static") || CheckNextToken((*reader),"field")){

	reader,wfile=ParseClassVarDec(reader,wfile)
	}
	
	for(CheckNextToken((*reader),"constructor") || CheckNextToken((*reader),"function") || CheckNextToken((*reader),"method")){

	reader,wfile=ParseSubDec(reader,wfile)
	}
	
	s,reader=getNextToken(reader) 
	
	wfile.WriteString(s+"\n") //}	
	wfile.WriteString("</class>")
	
	wfile.Close()
}

/***********ClassVarDec***********/

func ParseClassVarDec(reader *bufio.Reader ,wfile *os.File) (*bufio.Reader,*os.File) {
	s:=""
	if (CheckNextToken((*reader),"static") || CheckNextToken((*reader),"field")){

		wfile.WriteString("<classVarDec>\n")
		
		s,reader=getNextToken(reader) 
		wfile.WriteString(s+"\n")//declaretion
		
		s,reader=getNextToken(reader)
		wfile.WriteString(s+"\n")//type
    	
		s,reader=getNextToken(reader) 
		wfile.WriteString(s+"\n")//varname
		
		for (CheckNextToken((*reader),",")){
			
			s,reader=getNextToken(reader)
			wfile.WriteString(s+"\n")//,
			
			s,reader=getNextToken(reader) 
			wfile.WriteString(s+"\n")//varname
    	}
		s,reader=getNextToken(reader)
		wfile.WriteString(s+"\n")//;		
		}
		wfile.WriteString("</classVarDec>\n")
	 return reader,wfile
}


/***********subroutineDec***********/
func ParseSubDec(reader *bufio.Reader ,wfile *os.File) (*bufio.Reader,*os.File) {	
	s:=""
	if (CheckNextToken((*reader),"constructor") || CheckNextToken((*reader),"function") || CheckNextToken((*reader),"method")){
    	
		wfile.WriteString("<subroutineDec>\n")
		
			
			s,reader=getNextToken(reader) 
			wfile.WriteString(s+"\n")//type
			
			s,reader=getNextToken(reader)
			wfile.WriteString(s+"\n")//identifier
			
			s,reader=getNextToken(reader) 
			wfile.WriteString(s+"\n")//name
			
			s,reader=getNextToken(reader) 
			wfile.WriteString(s+"\n")//(
	
			reader,wfile=ParameterList(reader,wfile)
		
			s,reader=getNextToken(reader)
			wfile.WriteString(s+"\n")//)
			
			reader,wfile=SubRoutineBody(reader,wfile)
		
	
	wfile.WriteString("</subroutineDec>\n")
	}
	 return reader,wfile
}

/***********ParameterList***********/
func ParameterList(reader *bufio.Reader ,wfile *os.File) (*bufio.Reader,*os.File){
	s:=""
	wfile.WriteString("<parameterList>\n")
	if(!CheckNextToken((*reader),")")){
		
		s,reader=getNextToken(reader) 
		wfile.WriteString(s+"\n")//type
		
		s,reader=getNextToken(reader)
		wfile.WriteString(s+"\n")//varname
		for(CheckNextToken((*reader),",")){
		
			s,reader=getNextToken(reader) 
			wfile.WriteString(s+"\n")//,
			
			s,reader=getNextToken(reader) 
			wfile.WriteString(s+"\n")//type
			
			s,reader=getNextToken(reader) 
			wfile.WriteString(s+"\n")//varName
		}
		
	}
	wfile.WriteString("</parameterList>\n")

	return reader,wfile
}



/***********SubRoutineBody***********/
func SubRoutineBody(reader *bufio.Reader ,wfile *os.File) (*bufio.Reader,*os.File){
    	s:=""
		wfile.WriteString("<subroutineBody>\n")
		s,reader=getNextToken(reader) 
		wfile.WriteString(s+"\n")//{
		for(CheckNextToken((*reader),"var")){
		
		reader,wfile=varDec(reader,wfile)
		
		}
		statements(reader,wfile)
		
		s,reader=getNextToken(reader) 
		wfile.WriteString(s+"\n")//}
		
		wfile.WriteString("</subroutineBody>\n")
	
	return reader,wfile
}

/***********varDec***********/
func varDec(reader *bufio.Reader ,wfile *os.File) (*bufio.Reader,*os.File){
		s:=""
		wfile.WriteString("<varDec>\n")

		if(CheckNextToken((*reader),"var")){


				s,reader=getNextToken(reader) 
				wfile.WriteString(s+"\n")//var
				
				s,reader=getNextToken(reader) 
				wfile.WriteString(s+"\n")//type
				
				s,reader=getNextToken(reader) 
				wfile.WriteString(s+"\n")//varname
				
				for(CheckNextToken((*reader),",")){
					s,reader=getNextToken(reader) 
					wfile.WriteString(s+"\n")//,
					
					s,reader=getNextToken(reader) 
					wfile.WriteString(s+"\n")//varname
			}
			s,reader=getNextToken(reader) 
			wfile.WriteString(s+"\n")//;
	     	
}
	wfile.WriteString("</varDec>\n")

	return reader,wfile

}

/***********statments***********/
func statements(reader *bufio.Reader ,wfile *os.File) (*bufio.Reader,*os.File){
	
	wfile.WriteString("<statements>\n")
	for(CheckNextToken((*reader),"let")|| CheckNextToken((*reader),">if<") || CheckNextToken((*reader),"else") || CheckNextToken((*reader),"while") || CheckNextToken((*reader),"do") || CheckNextToken((*reader),"return")){
	   reader,wfile=statement(reader,wfile)
	}
	
	wfile.WriteString("</statements>\n")
	return reader,wfile
}


/***********statement***********/
func statement(reader *bufio.Reader ,wfile *os.File) (*bufio.Reader,*os.File){
    if(CheckNextToken((*reader),"let")){
	   reader,wfile=letStatement(reader,wfile)
	   }
	if(CheckNextToken((*reader),">if<")){
	  reader,wfile=ifStatement(reader,wfile)
	  }
	if(CheckNextToken((*reader),"while")){
		reader,wfile=whileStatement(reader,wfile)
		}
	if(CheckNextToken((*reader),"do")){
	reader,wfile=doStatement(reader,wfile)
	}
	if(CheckNextToken((*reader),"return")){
	reader,wfile=ReturnStatement(reader,wfile)
	}
	return reader,wfile
}


/***********letStatement***********/
func letStatement(reader *bufio.Reader ,wfile *os.File) (*bufio.Reader,*os.File) {	
	s:=""
	if(CheckNextToken((*reader),"let")){
		wfile.WriteString("<letStatement>\n")
		
		s,reader=getNextToken(reader)
		wfile.WriteString(s+"\n")//let
		
		s,reader=getNextToken(reader)
		wfile.WriteString(s+"\n")//varname
		
		if(CheckNextToken((*reader),"[")){
			s,reader=getNextToken(reader)
			wfile.WriteString(s+"\n")//[
		
			reader,wfile=expression(reader,wfile)

			s,reader=getNextToken(reader)
			wfile.WriteString(s+"\n")//]
		}
		s,reader=getNextToken(reader) 
		wfile.WriteString(s+"\n")//=
		
		reader,wfile=expression(reader,wfile)
	
	s,reader=getNextToken(reader)
	
	wfile.WriteString(s+"\n")//;	
	wfile.WriteString("</letStatement>\n")
	}
	 return reader,wfile
}


/***********ifStatement***********/
func ifStatement(reader *bufio.Reader ,wfile *os.File) (*bufio.Reader,*os.File){	
	s:=""
	if(CheckNextToken((*reader),">if<")){
		wfile.WriteString("<ifStatement>\n")
		
		s,reader=getNextToken(reader)
		wfile.WriteString(s+"\n")//if
		
		s,reader=getNextToken(reader)
		wfile.WriteString(s+"\n")//(
		
		reader,wfile=expression(reader,wfile)

		s,reader=getNextToken(reader)
		wfile.WriteString(s+"\n")//)
		
		s,reader=getNextToken(reader)
		wfile.WriteString(s+"\n")//{
		
		statements(reader,wfile)

		s,reader=getNextToken(reader)
		wfile.WriteString(s+"\n")//}
		if(CheckNextToken((*reader),"else")){
			
			s,reader=getNextToken(reader)
			wfile.WriteString(s+"\n")//else
			
			s,reader=getNextToken(reader)
			wfile.WriteString(s+"\n")//{
		
			reader,wfile=statements(reader,wfile)

			s,reader=getNextToken(reader)
			wfile.WriteString(s+"\n")//}
		}
			
		wfile.WriteString("</ifStatement>\n")
	}
	 return reader,wfile
}

/***********whileStatement***********/
func whileStatement(reader *bufio.Reader ,wfile *os.File) (*bufio.Reader,*os.File) {	
	s:=""
	if(CheckNextToken((*reader),"while")){
		
		wfile.WriteString("<whileStatement>\n")
		s,reader=getNextToken(reader)
		
		wfile.WriteString(s+"\n")//while
		s,reader=getNextToken(reader)
		
		wfile.WriteString(s+"\n")//(
		reader,wfile=expression(reader,wfile)

		s,reader=getNextToken(reader)
		wfile.WriteString(s+"\n")//)
		
		s,reader=getNextToken(reader)
		wfile.WriteString(s+"\n")//{
		
		statements(reader,wfile)

		s,reader=getNextToken(reader)
		wfile.WriteString(s+"\n")//}
		
		wfile.WriteString("</whileStatement>\n")
	}
	 return reader,wfile
}


/***********doStatement***********/
func doStatement(reader *bufio.Reader ,wfile *os.File) (*bufio.Reader,*os.File) {	
	s:=""
	if(CheckNextToken((*reader),"do")){
		wfile.WriteString("<doStatement>\n")
		s,reader=getNextToken(reader)
		wfile.WriteString(s+"\n")//do
	
		reader,wfile=subroutineCall(reader,wfile)

		s,reader=getNextToken(reader)
		wfile.WriteString(s+"\n")//;
		
		wfile.WriteString("</doStatement>\n")
	}
	 return reader,wfile
}


/***********ReturnStatement***********/
func ReturnStatement(reader *bufio.Reader ,wfile *os.File) (*bufio.Reader,*os.File) {	
	s:=""
	if(CheckNextToken((*reader),"return")){
		wfile.WriteString("<returnStatement>\n")
		
		s,reader=getNextToken(reader)
		wfile.WriteString(s+"\n")//return
	
		if(CheckNextToken((*reader),"integerConstant")||CheckNextToken((*reader),"stringConstant")||CheckNextToken((*reader),"true")||CheckNextToken((*reader),"false")||CheckNextToken((*reader),"null")||CheckNextToken((*reader),"this")||CheckNextToken((*reader),"identifier")||CheckNextToken((*reader),"(")||CheckNextToken((*reader),"-")||CheckNextToken((*reader),"~")){//if expression	
		    reader,wfile=expression(reader,wfile)
		}

		s,reader=getNextToken(reader)
		wfile.WriteString(s+"\n")//;
		
		wfile.WriteString("</returnStatement>\n")
	}
	 return reader,wfile
}


/***********expression***********/
func expression(reader *bufio.Reader ,wfile *os.File) (*bufio.Reader,*os.File) {	
		s:=""
		wfile.WriteString("<expression>\n")
		reader,wfile=term(reader,wfile)//term
		for(CheckNextToken((*reader),"+")||CheckNextToken((*reader),"*")||CheckNextToken((*reader),"-")||CheckNextToken((*reader),">/<")||CheckNextToken((*reader),"&amp;")||CheckNextToken((*reader),"|")||CheckNextToken((*reader),"&lt;")||CheckNextToken((*reader),"&gt;")||CheckNextToken((*reader),"=")){
			
			s,reader=getNextToken(reader)
			wfile.WriteString(s+"\n")//op
			
			reader,wfile=term(reader,wfile)
		}
		
		wfile.WriteString("</expression>\n")
	
	 return reader,wfile
}


/***********term***********/
func term(reader *bufio.Reader ,wfile *os.File) (*bufio.Reader,*os.File) {
		s:=""
		elseif:=true//internal
		elseif2:=true//external
		wfile.WriteString("<term>\n")
		
		if(CheckNextToken((*reader),"integerConstant")){
		
			s,reader=getNextToken(reader)
			wfile.WriteString(s+"\n")//integerConstant
			elseif2=false
	   }
	if(elseif2&&CheckNextToken((*reader),"stringConstant")){
	
			s,reader=getNextToken(reader)
			wfile.WriteString(s+"\n")//stringConstant
			elseif2=false
	   }
	if(elseif2&&CheckNextToken((*reader),"true")||CheckNextToken((*reader),"false")||CheckNextToken((*reader),"null")||CheckNextToken((*reader),"this")){
			
			s,reader=getNextToken(reader)
			wfile.WriteString(s+"\n")//stringConstant
			elseif2=false
	   }
	   
	if(elseif2&&CheckNextToken((*reader),"identifier")){
		if(CheckNextToken2((*reader),"[")){//var[
		
			s,reader=getNextToken(reader)
			wfile.WriteString(s+"\n")//varName
			
			s,reader=getNextToken(reader)
			wfile.WriteString(s+"\n")//[
		
			reader,wfile=expression(reader,wfile)

			s,reader=getNextToken(reader)
			wfile.WriteString(s+"\n")//]
			elseif=false
		}
		if(elseif && CheckNextToken2((*reader),"(")||CheckNextToken2((*reader),".")){
		//subcall
			reader,wfile=subroutineCall(reader,wfile)
			elseif=false
		}
		if(elseif){
			s,reader=getNextToken(reader)
			wfile.WriteString(s+"\n")//varName
		  elseif2=false   
		 }
	}
	if(elseif2&&CheckNextToken((*reader),"(")){
	
			s,reader=getNextToken(reader)
			wfile.WriteString(s+"\n")//(
			
			reader,wfile=expression(reader,wfile)
			
			s,reader=getNextToken(reader)
			wfile.WriteString(s+"\n")//)
			elseif2=false
	   }
	   
	if(elseif2&&CheckNextToken((*reader),"-")||CheckNextToken((*reader),"~")){
			
		s,reader=getNextToken(reader)
		wfile.WriteString(s+"\n")//unaryOp
			
		reader,wfile=term(reader,wfile)
	   }
	
		wfile.WriteString("</term>\n")
	
	 return reader,wfile
}


/***********subroutineCall***********/
func subroutineCall(reader *bufio.Reader ,wfile *os.File) (*bufio.Reader,*os.File) {	
		s:=""
	
		s,reader=getNextToken(reader)
		wfile.WriteString(s+"\n")//identifier
		
		if(CheckNextToken((*reader),"(")){
		    
			s,reader=getNextToken(reader)
			wfile.WriteString(s+"\n")//(
			
			reader,wfile=expressionList(reader,wfile)
			
			s,reader=getNextToken(reader)
			wfile.WriteString(s+"\n")//)
			
	   }
	if(CheckNextToken((*reader),".")){
			
			s,reader=getNextToken(reader)
			wfile.WriteString(s+"\n")//.
			
			s,reader=getNextToken(reader)
			wfile.WriteString(s+"\n")//subroutineName
			
			s,reader=getNextToken(reader)
			wfile.WriteString(s+"\n")//(
			
			reader,wfile=expressionList(reader,wfile)
			
			s,reader=getNextToken(reader)
			wfile.WriteString(s+"\n")//)
	   }
	
	 return reader,wfile
}

/***********expressionList***********/
func expressionList(reader *bufio.Reader ,wfile *os.File) (*bufio.Reader,*os.File) {	
		s:=""
		wfile.WriteString("<expressionList>\n")
		if(CheckNextToken((*reader),"integerConstant")||CheckNextToken((*reader),"stringConstant")||CheckNextToken((*reader),"true")||CheckNextToken((*reader),"false")||CheckNextToken((*reader),"null")||CheckNextToken((*reader),"this")||CheckNextToken((*reader),"identifier")||CheckNextToken((*reader),"(")||CheckNextToken((*reader),"-")||CheckNextToken((*reader),"~")){//if expression	
	
		   reader,wfile=expression(reader,wfile)
			for(CheckNextToken((*reader),",")){
					
					s,reader=getNextToken(reader) 
					wfile.WriteString(s+"\n")//,
					
					reader,wfile=expression(reader,wfile)
			}  
		}
		
		wfile.WriteString("</expressionList>\n")

	 return reader,wfile
}