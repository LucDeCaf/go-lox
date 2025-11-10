package ast

var Keywords map[string]TokenType

func init() {
	Keywords = make(map[string]TokenType)
	Keywords["and"] = AND
	Keywords["class"] = CLASS
	Keywords["else"] = ELSE
	Keywords["false"] = FALSE
	Keywords["for"] = FOR
	Keywords["fun"] = FUN
	Keywords["if"] = IF
	Keywords["nil"] = NIL
	Keywords["or"] = OR
	Keywords["print"] = PRINT
	Keywords["return"] = RETURN
	Keywords["super"] = SUPER
	Keywords["this"] = THIS
	Keywords["true"] = TRUE
	Keywords["var"] = VAR
	Keywords["while"] = WHILE
}
