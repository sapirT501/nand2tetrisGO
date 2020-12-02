package targil5
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

)
var if_index int
var while_index int
var classN string=""
var arg int =0//count var declarations
var fieldnum int=0//count field declarations
var staticnum int=0//count static declarations
var in int=0
var symbolClass [50]Table
type Table struct {  
    name string
    tipus string//int
    field string//field static
    index int
}
func Main(text string){

text = strings.TrimSuffix(text, "\n") //trimming not needed suffix, so file path will be read correctly
text = strings.TrimSuffix(text, "\r") 

str := text
str = filepath.ToSlash(str)

dir, err := ioutil.ReadDir(str)
if err != nil{
    log.Fatal(err)
  }

 
str= str + "/" //adding a slash to the end so it could lead to the file 


for _, f :=range dir { //loop over directory (the first part _ not needed because it is the index)
    fileName:=f.Name()
	if!f.IsDir() && strings.HasSuffix(fileName,"T.xml"){ //checking that f is not a directory and then making sure that the type of the file is "VM"
	    if_index=-1
		while_index=-1
		fieldnum=0
		in=0
		
		file,err := os.Open(str+f.Name())// open the TOKEN file for reading
		reader := bufio.NewReaderSize(file,50000) // getting a file reader
				
		fileNoExt := strings.TrimSuffix(f.Name(),"T.xml")//we dont want the type of file because according to the hack langauage it is not needed 
		
		/**creating xml file**/
		
		vmName := fileNoExt +".vm"
		file, err = os.Create(str + vmName)// Create- creates a file for reading and writing
        if err != nil{
           log.Fatal(err)
           }
		wfile, _ := os.OpenFile(str+vmName, os.O_WRONLY|os.O_APPEND, os.ModePerm)// we want to open the file for appending so if we write new lines the previous won't be deleted

		/***/
		
		if err == io.EOF {// making sure we are not at the end of the file
				break
		}
        ParseClass(reader ,wfile)
        




}//end of if

}//end of for dir


}//end of main


/****getNextToken****/
func getNextToken(reader *bufio.Reader) (string,*bufio.Reader){
	line, _, _ := reader.ReadLine()//reads the file line by line 
	s := string(line)// ReadeLine returns []byte but we want the line format to be a string
	return s,reader
	
}

/***func CheckNextToken***/
func CheckNextToken (reader bufio.Reader,s1 string) bool{
	line, _, _ := reader.ReadLine()//reads line by line the file. line is []byte
	s:= string(line)
	s= strings.TrimSuffix(s, "\n") //trimming not needed suffix, so filepath will be read correctly
	s= strings.TrimSuffix(s, "\r")
	if(strings.Contains(s,s1)){
		return true	
	}
	return false
}

/***func CheckNextToken2***/
func CheckNextToken2 (reader bufio.Reader,s1 string) bool{
	line, _, _ := reader.ReadLine()//reads  the file line by line line is [] byte
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
	s:=""
	
	_,reader=getNextToken(reader) //token of clasKeyword
	
	s,reader=getNextToken(reader) //token of className
	s=strings.TrimPrefix(s, "<identifier>")
	s=strings.TrimSuffix(s, "</identifier>")
	classN=s
	
	_,reader=getNextToken(reader)//{
	
	
	for(CheckNextToken((*reader),"static") || CheckNextToken((*reader),"field")){

	reader,wfile=ParseClassVarDec(reader,wfile)/// this is for the symbole table VM DONT CARE
	}
	
	for(CheckNextToken((*reader),"constructor") || CheckNextToken((*reader),"function") || CheckNextToken((*reader),"method")){

	reader,wfile=ParseSubDec(symbolClass,fieldnum,reader,wfile)///this is for the symbole table VM DONT CARE
	}
	
	_,reader=getNextToken(reader) //}
	
	wfile.Close()
}

/***********ClassVarDec***********/
//DONEEEEEEE//
func ParseClassVarDec( reader *bufio.Reader ,wfile *os.File) (*bufio.Reader,*os.File) {
	s:=""
	tipus:=""
	declar:=""
	varN:=""
	
	/*****====build the Symbole Table and then add the declarations====***/
	
    num:=0
	if (CheckNextToken((*reader),"static") || CheckNextToken((*reader),"field")){
		
		s,reader=getNextToken(reader) //declaretion field/static
		s=strings.TrimPrefix(s, "<keyword>")
		s=strings.TrimSuffix(s, "</keyword>")
		declar=s
		
		s,reader=getNextToken(reader)//type
		s=strings.TrimPrefix(s, "<keyword>")
		s=strings.TrimSuffix(s, "</keyword>")
		s=strings.TrimPrefix(s, "<identifier>")//identifier
		s=strings.TrimSuffix(s, "</identifier>")//identifier
		tipus=s
 
		s,reader=getNextToken(reader) //varname
		s=strings.TrimPrefix(s, "<identifier>")
		s=strings.TrimSuffix(s, "</identifier>")
		varN=s
		
		if(declar=="field"){declar="this"
		num=fieldnum
		fieldnum++
		}else{num=staticnum
		staticnum++}
		
		t:=Table{varN,tipus,declar,num}
		symbolClass[in]=t
		in++
		
		
		for (CheckNextToken((*reader),",")){
			
			s,reader=getNextToken(reader)//,
			
			
		    s,reader=getNextToken(reader) //varname
		    s=strings.TrimPrefix(s, "<identifier>")
		    s=strings.TrimSuffix(s, "</identifier>")
		    varN=s
			
			
			if(declar=="this"){
		        num=fieldnum
		        fieldnum++
		    }
			if(declar=="static"){num=staticnum
		          staticnum++}
				  
			//now lets build a table type//
			t=Table{varN,tipus,declar,in}
		    symbolClass[in]=t//add to symbol table
			in++
			
    	}
		s,reader=getNextToken(reader)//;
			
		}
	
	 return reader,wfile
}


/***********subroutineDec***********/
//DONEEEEEEE//
func ParseSubDec(symbolClass [50]Table,fieldN int,reader *bufio.Reader ,wfile *os.File) (*bufio.Reader,*os.File) {	
	s:=""
	index:=0
	tipus:=""
	methodN:=""
	i:=0
	var symbolMethod *[20]Table
	symbolMethod= new([20]Table)
	if (CheckNextToken((*reader),"constructor") || CheckNextToken((*reader),"function") || CheckNextToken((*reader),"method")){
    	
			s,reader=getNextToken(reader)//identifier =function
		    s=strings.TrimPrefix(s, "<keyword>")
		    s=strings.TrimSuffix(s, "</keyword>")
		    tipus=s
			
			s,reader=getNextToken(reader)// return type=void
	
		    s=strings.TrimPrefix(s, "<keyword>")
		    s=strings.TrimSuffix(s, "</keyword>")
			s=strings.TrimPrefix(s, "<identifier>")
		    s=strings.TrimSuffix(s, "</identifier>")
		    //declar=s
			
			s,reader=getNextToken(reader) //function name=main
			s=strings.TrimPrefix(s, "<identifier>")
		    s=strings.TrimSuffix(s, "</identifier>")
		    methodN=s

			 if(tipus=="method"){index=1}
			_,reader=getNextToken(reader) //(
			reader,wfile=ParameterList(symbolMethod,&index,&i,reader,wfile)//send the table to collect parameter
		
			_,reader=getNextToken(reader)//)
			 index=0
	         wfile.WriteString("function"+" "+classN+"."+methodN+" ")
			reader,wfile=SubRoutineBody(symbolClass,tipus,fieldN,symbolMethod,&index,&i,reader,wfile)
	        
	}
	 return reader,wfile
}

/***********ParameterList***********/
//DONEEEEEEE//
func ParameterList(symbolMethod *[20]Table,index *int,i *int,reader *bufio.Reader ,wfile *os.File) (*bufio.Reader,*os.File){
	s:=""
	tipus:=""
	declar:="argument"
	varN:=""
	
	if(!CheckNextToken((*reader),")")){
		
		s,reader=getNextToken(reader) //type
		s=strings.TrimPrefix(s, "<keyword>")
		s=strings.TrimSuffix(s, "</keyword>")
		s=strings.TrimPrefix(s, "<identifier>")
		s=strings.TrimSuffix(s, "</identifier>")
		tipus=s
		
		s,reader=getNextToken(reader)//varName
		s=strings.TrimPrefix(s, "<identifier>")
		s=strings.TrimSuffix(s, "</identifier>")
		varN=s
		
		t:=Table{varN,tipus,declar,*index}
		(*index)++
		(*symbolMethod)[*i]=t
		(*i)++
		
		for(CheckNextToken((*reader),",")){
		
			s,reader=getNextToken(reader)//,
			
			s,reader=getNextToken(reader)//type
			s=strings.TrimPrefix(s, "<keyword>")
			s=strings.TrimSuffix(s, "</keyword>")
			tipus=s			
			
			s,reader=getNextToken(reader)//varName
			s=strings.TrimPrefix(s, "<identifier>")
			s=strings.TrimSuffix(s, "</identifier>")
			varN=s
			
			t=Table{varN,tipus,declar,*index}
			(*index)++
			(*symbolMethod)[*i]=t
			(*i)++
		}
		
	}

	return reader,wfile
}



/***********SubRoutineBody***********/
//DONEEEEEEE//
func SubRoutineBody(symbolClass [50]Table,tipus string,fieldN int,symbolMethod *[20]Table,index *int,i *int,reader *bufio.Reader ,wfile *os.File) (*bufio.Reader,*os.File){

		_,reader=getNextToken(reader) //{

		for(CheckNextToken((*reader),"var")){
		
		reader,wfile=varDec(symbolMethod,index,i,reader,wfile)
		
		}
	    wfile.WriteString(strconv.Itoa(arg)+"\n")
		arg=0
        if(tipus=="method"){
		  	wfile.WriteString("push argument 0\npop pointer 0\n")       
		}
		 if(tipus=="constructor"){
		  	wfile.WriteString("push constant "+strconv.Itoa(fieldN)+"\ncall Memory.alloc 1\npop pointer 0\n")
		}
		statements(symbolClass,symbolMethod,reader,wfile)
		
		_,reader=getNextToken(reader) //}
			
	return reader,wfile
}

/***********varDec***********/

func varDec(symbolMethod *[20]Table,index *int,i *int,reader *bufio.Reader ,wfile *os.File) (*bufio.Reader,*os.File){
		s:=""
		tipus:=""
		declar:="local"
		varN:=""
		

		if(CheckNextToken((*reader),"var")){

			s,reader=getNextToken(reader) //var
			s=strings.TrimPrefix(s, "<keyword>")
			s=strings.TrimSuffix(s, "</keyword>")
			//declar=s	
				
			s,reader=getNextToken(reader) //type
			s=strings.TrimPrefix(s, "<keyword>")
			s=strings.TrimSuffix(s, "</keyword>")
			s=strings.TrimPrefix(s, "<identifier>")
			s=strings.TrimSuffix(s, "</identifier>")
			tipus=s	
				
			s,reader=getNextToken(reader) //varname
			s=strings.TrimPrefix(s, "<identifier>")
			s=strings.TrimSuffix(s, "</identifier>")
			varN=s
			
			t:=Table{varN,tipus,declar,*index}
			(*index)++
			symbolMethod[*i]=t
			(*i)++	
			arg++
			
			for(CheckNextToken((*reader),",")){
				s,reader=getNextToken(reader) //,
					
				s,reader=getNextToken(reader) //varname
				s=strings.TrimPrefix(s, "<identifier>")
				s=strings.TrimSuffix(s, "</identifier>")
				varN=s
			
				t:=Table{varN,tipus,declar,*index}
				(*index)++
				symbolMethod[*i]=t
				(*i)++	
				arg++
				
			}
			s,reader=getNextToken(reader) //;
	     	
}

	return reader,wfile

}

/***********statments***********/

func statements(symbolClass [50]Table,symbolMethod *[20]Table,reader *bufio.Reader ,wfile *os.File) (*bufio.Reader,*os.File){
	
	for(CheckNextToken((*reader),"let")|| CheckNextToken((*reader),">if<") || CheckNextToken((*reader),"else") || CheckNextToken((*reader),"while") || CheckNextToken((*reader),"do") || CheckNextToken((*reader),"return")){
	   reader,wfile=statement(symbolClass,symbolMethod,reader,wfile)
	}
	
	return reader,wfile
}


/***********statement***********/

func statement(symbolClass [50]Table,symbolMethod *[20]Table,reader *bufio.Reader ,wfile *os.File) (*bufio.Reader,*os.File){
    if(CheckNextToken((*reader),"let")){
	   reader,wfile=letStatement(symbolClass,symbolMethod,reader,wfile)
	   }
	if(CheckNextToken((*reader),">if<")){
	  reader,wfile=ifStatement(symbolClass,symbolMethod,reader,wfile)
	  }
	if(CheckNextToken((*reader),"while")){
		reader,wfile=whileStatement(symbolClass,symbolMethod,reader,wfile)
		}
	if(CheckNextToken((*reader),"do")){
	reader,wfile=doStatement(symbolClass,symbolMethod,reader,wfile)
	wfile.WriteString("pop temp 0\n")//pop the value returns from the called function to temp 0
	}
	if(CheckNextToken((*reader),"return")){
	reader,wfile=ReturnStatement(symbolClass,symbolMethod,reader,wfile)
	}
	return reader,wfile
}


/***********letStatement***********/
func letStatement(symbolClass [50]Table,symbolMethod *[20]Table,reader *bufio.Reader ,wfile *os.File) (*bufio.Reader,*os.File) {	
	s:=""
	if(CheckNextToken((*reader),"let")){
		
		_,reader=getNextToken(reader)//let
		
		s,reader=getNextToken(reader)//varname
		s=strings.TrimPrefix(s, "<identifier>")
		s=strings.TrimSuffix(s, "</identifier>")
		
		varN:=s
		if(CheckNextToken((*reader),"[")){
		        _,reader=getNextToken(reader)//[
		        //if in table***//
		        insymb:=false
				for _, b := range symbolMethod {
					if b.name == varN {
						insymb=true
						wfile.WriteString("push "+b.field +" "+strconv.Itoa(b.index)+"\n")
						break
					}
				}
			 if(!insymb){
			 
			  for _, b := range symbolClass {
				if b.name == varN {
				   wfile.WriteString("push "+b.field +" "+strconv.Itoa(b.index)+"\n")
				   break
                }
              }
			 
			 }
			reader,wfile=expression(symbolClass,symbolMethod,reader,wfile)
            wfile.WriteString("add\n")
			_,reader=getNextToken(reader)//]
			_,reader=getNextToken(reader) //=
		    reader,wfile=expression(symbolClass,symbolMethod, reader,wfile)
		    wfile.WriteString("pop temp 0\npop pointer 1\npush temp 0\npop that 0\n")
		}else{
		
		_,reader=getNextToken(reader) //=
		
		reader,wfile=expression(symbolClass,symbolMethod, reader,wfile)
		//check symbol table and then pop variable.declare variable.index
				        //if in table***//
		        insymb:=false
				for _, b := range symbolMethod {
					if b.name == varN {
						insymb=true
						wfile.WriteString("pop "+b.field +" "+strconv.Itoa(b.index)+"\n")
						break
					}
				}
			 if(!insymb){
			 
			  for _, b := range symbolClass {
               if b.name == varN {
				   wfile.WriteString("pop "+b.field +" "+strconv.Itoa(b.index)+"\n")
				   break
                }
              }
			 
			 }
			 
		}
	
	_,reader=getNextToken(reader)//;	
	}
	 return reader,wfile
}


/***********ifStatement***********/

func ifStatement(symbolClass [50]Table,symbolMethod *[20]Table,reader *bufio.Reader ,wfile *os.File) (*bufio.Reader,*os.File){	
	
	if(CheckNextToken((*reader),">if<")){
	    if_index++
		count:=if_index
		_,reader=getNextToken(reader)//if
		
		_,reader=getNextToken(reader)//(
		
		reader,wfile=expression(symbolClass,symbolMethod,reader,wfile)
	
		wfile.WriteString("if-goto IF_TRUE"+strconv.Itoa(count)+"\ngoto IF_FALSE"+strconv.Itoa(count)+"\nlabel IF_TRUE"+strconv.Itoa(count)+"\n")

		_,reader=getNextToken(reader)//)
		
		_,reader=getNextToken(reader)//{
		
		statements(symbolClass,symbolMethod,reader,wfile)
		
		wfile.WriteString("goto IF_END"+strconv.Itoa(count)+"\nlabel IF_FALSE"+strconv.Itoa(count)+"\n")

		_,reader=getNextToken(reader)//}

		if(CheckNextToken((*reader),"else")){
			
			_,reader=getNextToken(reader)//else
			
			_,reader=getNextToken(reader)//{
		   
			reader,wfile=statements(symbolClass,symbolMethod,reader,wfile)
			

			_,reader=getNextToken(reader)//}
		}
		wfile.WriteString("label IF_END"+strconv.Itoa(count)+"\n")
		
	}
	 return reader,wfile
}

/***********whileStatement***********/

func whileStatement(symbolClass [50]Table,symbolMethod *[20]Table,reader *bufio.Reader ,wfile *os.File) (*bufio.Reader,*os.File) {	
	
	if(CheckNextToken((*reader),"while")){
		while_index++
		count:=while_index
		
		_,reader=getNextToken(reader)//while
		
		_,reader=getNextToken(reader)//(
		
		wfile.WriteString("label WHILE_EXP"+strconv.Itoa(count)+"\n")

		reader,wfile=expression(symbolClass,symbolMethod,reader,wfile)

		wfile.WriteString("not\nif-goto WHILE_END"+strconv.Itoa(count)+"\n")

		_,reader=getNextToken(reader)//)
		
		_,reader=getNextToken(reader)//{
		
		statements(symbolClass,symbolMethod,reader,wfile)
		wfile.WriteString("goto WHILE_EXP"+strconv.Itoa(count)+"\nlabel WHILE_END"+strconv.Itoa(count)+"\n")
		
		_,reader=getNextToken(reader)//}
		
		
	}
	 return reader,wfile
}


/***********doStatement***********/

func doStatement(symbolClass [50]Table,symbolMethod *[20]Table,reader *bufio.Reader ,wfile *os.File) (*bufio.Reader,*os.File) {	

	if(CheckNextToken((*reader),"do")){
		_,reader=getNextToken(reader)//do
	     
		reader,wfile=subroutineCall(symbolClass,symbolMethod,reader,wfile)

		_,reader=getNextToken(reader)//;
	}
return reader,wfile
}


/***********ReturnStatement***********/
func ReturnStatement(symbolClass [50]Table,symbolMethod *[20]Table,reader *bufio.Reader ,wfile *os.File) (*bufio.Reader,*os.File) {	
	
	if(CheckNextToken((*reader),"return")){
		
		_,reader=getNextToken(reader)//return
	
		if(CheckNextToken((*reader),"integerConstant")||CheckNextToken((*reader),"stringConstant")||CheckNextToken((*reader),"true")||CheckNextToken((*reader),"false")||CheckNextToken((*reader),"null")||CheckNextToken((*reader),"this")||CheckNextToken((*reader),"identifier")||CheckNextToken((*reader),"(")||CheckNextToken((*reader),"-")||CheckNextToken((*reader),"~")){//if expression	
		    reader,wfile=expression(symbolClass,symbolMethod,reader,wfile)
		}else{
			wfile.WriteString("push constant 0\n")
	 }

		_,reader=getNextToken(reader)//;
		wfile.WriteString("return\n")

	}
	 return reader,wfile
}


/***********expression***********/


func expression(symbolClass [50]Table,symbolMethod *[20]Table,reader *bufio.Reader ,wfile *os.File) (*bufio.Reader,*os.File) {	

		reader,wfile=term(symbolClass,symbolMethod,reader,wfile)//term
		for(CheckNextToken((*reader),"+")||CheckNextToken((*reader),"*")||CheckNextToken((*reader),"-")||CheckNextToken((*reader),">/<")||CheckNextToken((*reader),"&amp;")||CheckNextToken((*reader),"|")||CheckNextToken((*reader),"&lt;")||CheckNextToken((*reader),"&gt;")||CheckNextToken((*reader),"=")){
			
			//s,reader=getNextToken(reader)

			if(CheckNextToken((*reader),"+")){
			    _,reader=getNextToken(reader)
			    reader,wfile=term(symbolClass,symbolMethod,reader,wfile)
				wfile.WriteString("add\n")
			}
			
			if(CheckNextToken((*reader),"-")){
				_,reader=getNextToken(reader)
				reader,wfile=term(symbolClass,symbolMethod,reader,wfile)
				wfile.WriteString("sub\n")
			}
			
			if(CheckNextToken((*reader),"*")){
				_,reader=getNextToken(reader)
				reader,wfile=term(symbolClass,symbolMethod,reader,wfile)
				wfile.WriteString("call Math.multiply 2\n")
			}
			
			if(CheckNextToken((*reader),">/<")){
				_,reader=getNextToken(reader)
				reader,wfile=term(symbolClass,symbolMethod,reader,wfile)
				wfile.WriteString("call Math.divide 2\n")
			}
			
			if(CheckNextToken((*reader),"&amp;")){
				_,reader=getNextToken(reader)
				reader,wfile=term(symbolClass,symbolMethod,reader,wfile)
				wfile.WriteString("and\n")
			}
			
			if(CheckNextToken((*reader),"|")){
				_,reader=getNextToken(reader)
				reader,wfile=term(symbolClass,symbolMethod,reader,wfile)
				wfile.WriteString("or\n")
			}
			
			if(CheckNextToken((*reader),"&lt;")){
				_,reader=getNextToken(reader)
				reader,wfile=term(symbolClass,symbolMethod,reader,wfile)
				wfile.WriteString("lt\n")
			}
			
			if(CheckNextToken((*reader),"&gt;")){
				_,reader=getNextToken(reader)
				reader,wfile=term(symbolClass,symbolMethod,reader,wfile)
				wfile.WriteString("gt\n")
			}
			
			if(CheckNextToken((*reader),"=")){
				_,reader=getNextToken(reader)
				reader,wfile=term(symbolClass,symbolMethod,reader,wfile)
				wfile.WriteString("eq\n")
			}
		}
		
	
	 return reader,wfile
}


/***********term***********/
func term(symbolClass [50]Table,symbolMethod *[20]Table,reader *bufio.Reader ,wfile *os.File) (*bufio.Reader,*os.File) {
		s:=""
		elseif:=true//internal
		elseif2:=true//external
	
		
		if(CheckNextToken((*reader),"integerConstant")){

			s,reader=getNextToken(reader)
			s=strings.TrimPrefix(s, "<integerConstant>")
			s=strings.TrimSuffix(s, "</integerConstant>")
			wfile.WriteString("push constant "+s+"\n")//integerConstant
			elseif2=false
	   }
	if(elseif2&&CheckNextToken((*reader),"stringConstant")){
	
			s,reader=getNextToken(reader)
			s=strings.TrimPrefix(s, "<stringConstant>")
			s=strings.TrimSuffix(s, "</stringConstant>")
			wfile.WriteString("push constant ")
			wfile.WriteString(strconv.Itoa(len(s)))
			wfile.WriteString("\ncall String.new 1\n")//stringConstant
			for _, c :=range s  {          
	           temp:=int(c)
	         wfile.WriteString("push constant ")
			 wfile.WriteString(strconv.Itoa(temp))
			 wfile.WriteString("\ncall String.appendChar 2\n")
			 }
			elseif2=false

	   }
	   
	if(elseif2&&CheckNextToken((*reader),"true")||CheckNextToken((*reader),"false")||CheckNextToken((*reader),"null")||CheckNextToken((*reader),"this")){
			
			
			if(CheckNextToken((*reader),"true")){
			    _,reader=getNextToken(reader)
				wfile.WriteString("push constant 0\nnot\n")
			}
			if(CheckNextToken((*reader),"false")||CheckNextToken((*reader),"null")){
				_,reader=getNextToken(reader)
				wfile.WriteString("push constant 0\n")
			}
			
			if(CheckNextToken((*reader),"this")){
				_,reader=getNextToken(reader)
				wfile.WriteString("push pointer 0\n")
			}
			elseif2=false
	   }
	   
	if(elseif2&&CheckNextToken((*reader),"identifier")){
		if(CheckNextToken2((*reader),"[")){//var[
		
			s,reader=getNextToken(reader)//varName
			s=strings.TrimPrefix(s, "<identifier>")
			s=strings.TrimSuffix(s, "</identifier>")
			varN:=s
			_,reader=getNextToken(reader)//[
		
			reader,wfile=expression(symbolClass,symbolMethod,reader,wfile)

			_,reader=getNextToken(reader)//]
			insymb:=false
			for _, b := range symbolMethod {
               if b.name == varN {
                   insymb=true
				   wfile.WriteString("push "+b.field +" "+strconv.Itoa(b.index)+"\n")
				   break
               }
             }
			 if(!insymb){
			 
			  for _, b := range symbolClass {
               if b.name == varN {
				   wfile.WriteString("push "+b.field +" "+strconv.Itoa(b.index)+"\n")
				   break
                }
              }
			 
			 }
			wfile.WriteString("\nadd\npop pointer 1\npush that 0\n")
			elseif=false
		}
		if(elseif && CheckNextToken2((*reader),"(")||CheckNextToken2((*reader),".")){
		//subcall
			reader,wfile=subroutineCall(symbolClass,symbolMethod,reader,wfile)
			elseif=false
		}
		if(elseif){// single variable
		    insymb:=false
			s,reader=getNextToken(reader)
			s=strings.TrimPrefix(s, "<identifier>")
			s=strings.TrimSuffix(s, "</identifier>")
			for _, b := range symbolMethod {
               if b.name == s {
                   insymb=true
				   wfile.WriteString("push "+b.field +" "+strconv.Itoa(b.index)+"\n")
				   break
               }
             }
			 if(!insymb){
			 
			  for _, b := range symbolClass {
               if b.name == s {
				   wfile.WriteString("push "+b.field +" "+strconv.Itoa(b.index)+"\n")
				   break
                }
              }
			 
			 }
		  elseif2=false   
		 }
	}
	if(elseif2&&CheckNextToken((*reader),"(")){
	
			s,reader=getNextToken(reader)
		
			reader,wfile=expression(symbolClass,symbolMethod,reader,wfile)
			
			s,reader=getNextToken(reader)
			
			elseif2=false
	   }
	   
	if(elseif2&&CheckNextToken((*reader),"-")||CheckNextToken((*reader),"~")){
		
		if(CheckNextToken((*reader),"-")){
			s,reader=getNextToken(reader)
			reader,wfile=term(symbolClass,symbolMethod ,reader,wfile)
			wfile.WriteString("neg\n")
			}
			
		if(CheckNextToken((*reader),"~")){
			s,reader=getNextToken(reader)
			reader,wfile=term(symbolClass,symbolMethod ,reader,wfile)
			wfile.WriteString("not\n")
		}
		
		
	   }
		
	 return reader,wfile
}


/***********subroutineCall***********/

func subroutineCall(symbolClass [50]Table,symbolMethod *[20]Table,reader *bufio.Reader ,wfile *os.File) (*bufio.Reader,*os.File) {	
		s:=""
		param:=0
		param1:=0
		s,reader=getNextToken(reader)
		s=strings.TrimPrefix(s, "<identifier>")
	    s=strings.TrimSuffix(s, "</identifier>")
		funcN:=s//function Name or other class Name or object Name
		
		if(CheckNextToken((*reader),"(")){//my method
		    
			s,reader=getNextToken(reader)//(
			
			wfile.WriteString("push pointer 0\n")
			reader,wfile,param=expressionList(symbolClass,symbolMethod,reader,wfile)
			param++
			wfile.WriteString("call "+classN+"."+funcN+" "+strconv.Itoa(param)+"\n")
			s,reader=getNextToken(reader)//)
			
			
	   }
	if(CheckNextToken((*reader),".")){//function/another class
			
			s,reader=getNextToken(reader)//.
			
			s,reader=getNextToken(reader)//subroutineName
			s=strings.TrimPrefix(s, "<identifier>")
	        s=strings.TrimSuffix(s, "</identifier>")
			_,reader=getNextToken(reader)//(
			
			 insymb:=false
			 //checking if the identifier is an object
			 for _, b := range symbolMethod {
               if b.name == funcN {
                   insymb=true
				   wfile.WriteString("push "+b.field+" "+strconv.Itoa(b.index)+"\n")
		                   funcN=b.tipus//if it is object we want (the class name)
				   param1++
				   break
               }
             }
			 if(!insymb){
			 
			  for _, b := range symbolClass {
               if b.name == funcN {
				   wfile.WriteString("push "+b.field+" "+strconv.Itoa(b.index)+"\n")
				   funcN=b.tipus
				   param1++
				   break
                }
              }
			 
			 }
			reader,wfile,param=expressionList(symbolClass,symbolMethod,reader,wfile)
			param=param+param1
			wfile.WriteString("call "+funcN+"."+s+" "+strconv.Itoa(param)+"\n")
			
			s,reader=getNextToken(reader)//)
	   }
	
	 return reader,wfile
}

/***********expressionList***********/
//DONE///////////
func expressionList(symbolClass [50]Table,symbolMethod *[20]Table,reader *bufio.Reader ,wfile *os.File) (*bufio.Reader,*os.File,int) {	
        param:=0
		if(CheckNextToken((*reader),"integerConstant")||CheckNextToken((*reader),"stringConstant")||CheckNextToken((*reader),"true")||CheckNextToken((*reader),"false")||CheckNextToken((*reader),"null")||CheckNextToken((*reader),"this")||CheckNextToken((*reader),"identifier")||CheckNextToken((*reader),"(")||CheckNextToken((*reader),"-")||CheckNextToken((*reader),"~")){//if expression	
	      
		   reader,wfile=expression(symbolClass,symbolMethod ,reader,wfile)
		   param++
			for(CheckNextToken((*reader),",")){
					
					_,reader=getNextToken(reader)//,
					//wfile.WriteString("push ")
					reader,wfile=expression(symbolClass,symbolMethod,reader,wfile)
					param++
			}
			
		}
		
	

	 return reader,wfile,param
}
