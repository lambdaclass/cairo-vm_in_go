func main() {
    fib(1, 1, 10000);

    ret;
}

func fib(first_element, second_element, n) -> (res: felt) {
    if (n == 0) {
        return (second_element,);
    }

    tempvar y = first_element + second_element;
    return fib(second_element, y, n - 1);
}