package main

import (
	//	"context"
	"fmt"
	"math/rand"
	"time"
)

/* //БЛОК1
//В первом блоке небольшое обучение по каналам и горутинам

// так помечается канал только на чтение func writer() <-chan int

func writer() <-chan int {
	//канал обязательно нужно инициализировать через команду make (использовать var для этих целей нельзя)
	ch := make(chan int)
	wg := &sync.WaitGroup{}

	wg.Add(2)

	//если не использовать go рутину и писать в канал то будет ошибка fatal error: all goroutines are asleep - deadlock! goroutine 1 [chan send]
	go func() {
		defer wg.Done()
		for i := range 5 {
			ch <- i + 1
		}
		//мы закрываем канал как только перестаем писать в него. Это нужно для того чтобы в main было понятно что в канал уже ничего не пишут
		// но если мы используем синхронизацию wg &sync то закрывать канал нужно уже в другой рутине
		// close(ch)
	}()

	go func() {
		defer wg.Done()
		for i := range 5 {
			ch <- i + 11
		}
		//мы закрываем канал как только перестаем писать в него. Это нужно для того чтобы в main было понятно что в канал уже ничего не пишут
		// но если мы используем синхронизацию wg &sync то закрывать канал нужно уже в другой рутине
		// close(ch)
	}()

	go func() {
		wg.Wait()
		//		close(ch)
	}()
	return ch
}

func main() {
	ch := writer()

	// если закрыть канал токо для чтения будет ошибка компиляцииclose(ch)
	//если мы знаем количество считываний из канала то можно и не закрывать канал
	for v := range ch {
		// v := <-ch
		//код в комменте можно заменить чтение через range
		//		v, ok := <-ch
		//в ок записывается true false в зависимости от того есть ли v
		//		if !ok {
		//			break
		//		}
		fmt.Println("v =", v)
	}

	time.Sleep(1 * time.Second)
}*/

/*//БЛОК2
// небольшая программка которая в одной гоу рутине создает канал с числами от 1 - 10, потом во второй рутине
// умножает каждое сило на 2 и в осной ной программе выводит эти числа с небольшой задержкой
func reader(ch <-chan int) {

	for i := range ch {
		fmt.Println(i)
	}
}

func doubler(ch <-chan int) <-chan int {
	ch1 := make(chan int)

	go func() {
		for i := range ch {
			time.Sleep(1 * time.Second)
			ch1 <- i * 2
		}
		close(ch1)
	}()
	return ch1
}

func writer() <-chan int {
	ch := make(chan int)

	go func() {
		for i := range 10 {
			ch <- i + 1
		}
		close(ch)
	}()
	return ch
}

func main() {
	reader(doubler(writer()))
}
*/

/*//БЛОК3
//Теория по чтению и записи в один канал
func main() {
	ch := make(chan int)

	go func() {
		for i := range 100 {
			ch <- i
		}
		close(ch)

	}()

	go func() {
		for i := range 100 {
			ch <- i * 2
		}
		//close(ch)

	}()

	go func() {
		for i := range ch {
			fmt.Println("i=", i, "worker1")
		}
	}()

	for i := range ch {
		fmt.Println("i=", i, "Worker2")
	}
}
*/

/*//БЛОК4
//Использование select

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		ch1 <- 1 //ch2 <- 1
	}()

	timer := time.NewTimer(1 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond)
	defer cancel()

	ch3 := make(chan int)
	close(ch3)

	select { //это блокирующий оператор
	case v := <-ch1: //если с канал ничего не пришло или программа сработет раньше чем рутина то и читать нечего по этому сработет default
		fmt.Println("v=", v, "from ch1")
	case v := <-ch2: //по этому в рутине должна быть задержка
		fmt.Println("v=", v, "from ch2")
	case <-time.After(1 * time.Second): // еще вариант срабатывания сейса по функции After
		fmt.Println("exit by after")
	case <-timer.C: //или по таймеру
		fmt.Println("exit by timer")
	case <-ctx.Done(): //или по контексту
		fmt.Println("exit by context")
	default:
		fmt.Println("exit by default")
	}
}
*/

// БЛОК5 Задачи
// есть функция которая работает 100 секунд
func randomTimeWork() {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Second)
}

//нужно написать обертку для этой функции которая будет прерывать выполнение
// если функция работает больше 3 сек и возвращать ошибку

func predictableTimeWork() error {
	ch := make(chan struct{})
	var err error

	go func() {
		randomTimeWork()
		close(ch)
	}()

	select {
	case <-ch:
		fmt.Println("func exequt")
		return nil
	case <-time.After(3 * time.Second):
		fmt.Println("error timeout")
		return err
	}
}

func main() {
	predictableTimeWork()
}
