1. Написать функцию, которая принимает на вход структуру `in` (`struct` или кастомную `struct`) и
   `values map[string]interface{}` (`key` - название поля структуры, которому нужно присвоить value
   этой мапы). Необходимо по значениям из мапы изменить входящую структуру `in` с помощью
   пакета `reflect`. Функция может возвращать только ошибку `error`. Написать к данной функции
   тесты (чем больше, тем лучше - зачтется в плюс). 
   
   > cd lsn7
   
   > go test struct_loader/* -v


      === RUN   TestLoadSimpleStruct
      === RUN   TestLoadSimpleStruct/StructSimple1_convert
      === RUN   TestLoadSimpleStruct/StructSimple1
      === RUN   TestLoadSimpleStruct/StructSimple2
      === RUN   TestLoadSimpleStruct/StructSimple3
      --- PASS: TestLoadSimpleStruct (0.00s)
      --- PASS: TestLoadSimpleStruct/StructSimple1_convert (0.00s)
      --- PASS: TestLoadSimpleStruct/StructSimple1 (0.00s)
      --- PASS: TestLoadSimpleStruct/StructSimple2 (0.00s)
      --- PASS: TestLoadSimpleStruct/StructSimple3 (0.00s)
      === RUN   TestLoadSimpleStructErrors
      === RUN   TestLoadSimpleStructErrors/ErrorIsNil
      === RUN   TestLoadSimpleStructErrors/ErrorInvalidReceiver
      === RUN   TestLoadSimpleStructErrors/ErrorNotConvertable
      --- PASS: TestLoadSimpleStructErrors (0.00s)
      --- PASS: TestLoadSimpleStructErrors/ErrorIsNil (0.00s)
      --- PASS: TestLoadSimpleStructErrors/ErrorInvalidReceiver (0.00s)
      --- PASS: TestLoadSimpleStructErrors/ErrorNotConvertable (0.00s)
      === RUN   TestLoadStructNamed
      === RUN   TestLoadStructNamed/StructNamed1
      === RUN   TestLoadStructNamed/StructNamed2
      --- PASS: TestLoadStructNamed (0.00s)
      --- PASS: TestLoadStructNamed/StructNamed1 (0.00s)
      --- PASS: TestLoadStructNamed/StructNamed2 (0.00s)
      === RUN   TestLoadStructDefault
      === RUN   TestLoadStructDefault/StructDefault1
      --- PASS: TestLoadStructDefault (0.00s)
      --- PASS: TestLoadStructDefault/StructDefault1 (0.00s)
      === RUN   TestLoadStructRequiredErrors
      === RUN   TestLoadStructRequiredErrors/StructRequired_Bool
      === RUN   TestLoadStructRequiredErrors/StructRequired_Int
      === RUN   TestLoadStructRequiredErrors/StructRequired_Float
      === RUN   TestLoadStructRequiredErrors/StructRequired_TD
      === RUN   TestLoadStructRequiredErrors/StructRequired_Slice
      === RUN   TestLoadStructRequiredErrors/StructRequired_Map
      === RUN   TestLoadStructRequiredErrors/StructRequired_Struct
      --- PASS: TestLoadStructRequiredErrors (0.00s)
      --- PASS: TestLoadStructRequiredErrors/StructRequired_Bool (0.00s)
      --- PASS: TestLoadStructRequiredErrors/StructRequired_Int (0.00s)
      --- PASS: TestLoadStructRequiredErrors/StructRequired_Float (0.00s)
      --- PASS: TestLoadStructRequiredErrors/StructRequired_TD (0.00s)
      --- PASS: TestLoadStructRequiredErrors/StructRequired_Slice (0.00s)
      --- PASS: TestLoadStructRequiredErrors/StructRequired_Map (0.00s)
      --- PASS: TestLoadStructRequiredErrors/StructRequired_Struct (0.00s)
      === RUN   TestLoadStructNested
      === RUN   TestLoadStructNested/StructNested1
      === RUN   TestLoadStructNested/StructNested2
      === RUN   TestLoadStructNested/StructNested3
      --- PASS: TestLoadStructNested (0.00s)
      --- PASS: TestLoadStructNested/StructNested1 (0.00s)
      --- PASS: TestLoadStructNested/StructNested2 (0.00s)
      --- PASS: TestLoadStructNested/StructNested3 (0.00s)
      PASS
      ok      command-line-arguments  (cached)


2. Написать функцию, которая принимает на вход имя файла и название функции. Необходимо
подсчитать в этой функции количество вызовов асинхронных функций. Результат работы
должен возвращать количество вызовов `int` и ошибку `error`. Разрешается использовать только
`go/parser`, `go/ast` и `go/token`.
   > cd lsn7
   
   > go test -v goanalyzer/*
   
  
      === RUN   TestGetGoFuncCallCount
      === RUN   TestGetGoFuncCallCount/GetGoFuncCallCount4sieve_channel
      --- PASS: TestGetGoFuncCallCount (0.00s)
      --- PASS: TestGetGoFuncCallCount/GetGoFuncCallCount4sieve_channel (0.00s)
      PASS
      ok      command-line-arguments  0.003s
   

3. *не обязательное*. Написать кодогенератор под какую-нибудь задачу.
