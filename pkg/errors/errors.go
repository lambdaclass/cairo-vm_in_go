package main

import (
	"fmt"
	"os"

	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/pkg/errors"
)

var CustomErrorVar = errors.New("Defined Custom Error")

// CustomError is a custom error type.
type CustomError struct {
	err     error
	Message string
}

func (ce *CustomError) Error() string {
	return ce.Message
}

func main() {
	fmt.Println("Starting the program.")

	// Handling Multiple Errors
	err1 := errors.New("First error")
	err2 := errors.New("Second error")
	if err := handleMultipleErrors(err1, err2); err != nil {
		fmt.Println("Multiple Errors:", err)
	}

	// Using Wrapf and More
	useWrapf()
	useErrorf()
	useWrap()

	// Custom Error
	customErr := &CustomError{Message: "This is a custom error"}
	fmt.Println("\nCustom Error:", customErr)

	fmt.Println("End of the program.")

	// Defer and Errors
	file := createFile("test.txt")
	defer closeFile(file)
	writeToFile(file, "Hello, World!")

	// Panic and Recover
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()
	panic("A panic occurred!")
}

func handleMultipleErrors(errorsToHandle ...error) error {
	var allErrors []string

	for _, err := range errorsToHandle {
		if err != nil {
			allErrors = append(allErrors, err.Error())
		}
	}

	if len(allErrors) > 0 {
		return errors.New("Multiple errors occurred: " + fmt.Sprintf("%v", allErrors))
	}

	return nil
}

func createFile(filename string) *os.File {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return nil
	}
	return file
}

func writeToFile(file *os.File, text string) {
	if file != nil {
		_, err := file.WriteString(text)
		if err != nil {
			fmt.Println("Error writing to file:", err)
		}
	}
}

func closeFile(file *os.File) {
	if file != nil {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
	}
}

func useWrapf() {
	originalErr := errors.New("original error")
	wrappedErr := errors.Wrapf(originalErr, "additional context: %T", lambdaworks.FeltOne())

	fmt.Println("\nUsing errors.Wrapf():")
	fmt.Println("Original Error:", originalErr)
	fmt.Println("Wrapped Error:", wrappedErr)
}

func useErrorf() {
	wrappedErr := errors.Errorf("error with additional context: %s", "some context")

	fmt.Println("\nUsing errors.Errorf():")
	fmt.Println("Wrapped Error:", wrappedErr)
}

func useWrap() {
	originalErr := errors.New("original error")
	wrappedErr := errors.Wrap(originalErr, "additional context")

	fmt.Println("\nUsing errors.Wrap():")
	fmt.Println("Original Error:", originalErr)
	fmt.Println("Wrapped Error:", wrappedErr)
}
