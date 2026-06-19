package main

import (
	"hash/fnv"
	"sync"
)

// Дан интерфейс нужно показать работу с кэшем
type Cache interface {
	Set(k string, v string)
	Get(k string) (string, bool)
}

/* Эта структура при шардировании земеняется на структуру Shard
//Создаем структуру для реализации интерфейса. Для хранения данных будем использовать тип данных "map"
type InMemoryCache struct {
	data map[string]string
	mu   sync.RWMutex // а так же во второй части решаем проблему конкурентного достума при помощи RWMutex
}
// Эта структура при шардировании земеняется на структуру Shard
*/

// Новая структура InMemoryCache содержит слайс Shard ов
type InMemoryCache struct {
	shards []Shard
}

// Новая реализация NewInMemoryCache (где инициализируется каждая мапа)
func NewInMemoryCache(numShards int) *InMemoryCache {
	shards := make([]Shard, 0, numShards)
	for i := 0; i < numShards; i++ {
		shards = append(shards, Shard{data: make(map[string]string)})
	}

	return &InMemoryCache{
		shards: shards,
	}
}

// Тут к структуре InMemoryCache создаются два метода Set и Get которые шардируют данные по хэшу
func (c *InMemoryCache) Set(k string, v string) {
	shardID := hasher(k) % len(c.shards)
	println(shardID)
	c.shards[shardID].Set(k, v) //так как shards это элемент структуры InMemoryCache с типом данных структура []Shard, а у []Shard
	// есть сетоды Set и Get то мы можем их сдесь применять
}

// Тут к структуре InMemoryCache создаются два метода Set и Get которые шардируют данные по хэшу
func (c *InMemoryCache) Get(k string) (string, bool) {
	shardID := hasher(k) % len(c.shards)
	return c.shards[shardID].Get(k) //так как shards это элемент структуры InMemoryCache с типом данных структура []Shard, а у []Shard
	// есть сетоды Set и Get то мы можем их сдесь применять
}

// функция hasher рандомно генерит число в формате uint32 используя пакет fnv
func hasher(k string) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(k))
	return int(h.Sum32())

}

// Для того чтобы сделать шардирование структуру InMemoryCache переименовываем в Shsrd
type Shard struct {
	data map[string]string
	mu   sync.RWMutex
}

/* Эта структура при шардировании изменяется
// В этой функции мы реализуем инициализацию мапы с помощью функции "make". без этого мы не сможем ничего записать в Nilовую мапу
func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		data: make(map[string]string),
	}
}
// Эта структура при шардировании изменяется
*/

// создаем методы описанные в интерфейсе для (c *InMemoryCache) при шардировании мы меняем на (s *Shard)
func (s *Shard) Set(k string, v string) {
	s.mu.Lock()         //Вторая часть задачи с блокировкой доступа. Так мы блокируем доступ для записи к мапе
	defer s.mu.Unlock() // и обязательно сразу в дифере разблокируем доступ
	s.data[k] = v
}

// создаем методы описанные в интерфейсе для (c *InMemoryCache) при шардировании мы меняем на (s *Shard)
func (s *Shard) Get(k string) (string, bool) {
	s.mu.RLock()         //Вторая часть задачи с блокировкой доступа. Тут мы используем блок для записи во время чтения
	defer s.mu.RUnlock() // и обязательно сразу в дифере разблокируем доступ
	data, ok := s.data[k]
	return data, ok
}

func main() {
	/* Этот код реализует простую запись в кэш
	// Этот код реализует простую запись в кэш
	cache := NewInMemoryCache()
	cache.Set("foo", "bar")

	data, ok := cache.Get("foo")
	if !ok {
		println("Key foo not found")
		return
	}
	println("Key: foo Value:", data)
	// Конец кода реализии простой записи в кэш
	*/

	/*этот код уже работат с использование Mutex
	//тут запускается несколько запросов Set, а так же в Go рутинах запускаются процессы Set и Get имитируя конкурентный доступ
	cache := NewInMemoryCache(5)
	cache.Set("foo", "bar")
	cache.Set("boo", "bah")

	wg := sync.WaitGroup{} // так же добавлено использование WaitGroup для того чтобы обеспечить выполнение рутин до завершения main
	wg.Add(4)              //эта функция говорит о том сколько рутин будет запущено

	go func() {
		defer wg.Done() // в дифере запускается Done который говорит о том что рутина сработала и делает -1 от числа Add
		cache.Set("foo", "update_bar")
		println("1")
	}()

	go func() {
		defer wg.Done() // в дифере запускается Done который говорит о том что рутина сработала и делает -1 от числа Add
		cache.Set("boo", "update_bah")
		println("2")
	}()

	go func() {
		defer wg.Done() // в дифере запускается Done который говорит о том что рутина сработала и делает -1 от числа Add
		cache.Get("foo")
		println("3")
	}()

	go func() {
		defer wg.Done() // в дифере запускается Done который говорит о том что рутина сработала и делает -1 от числа Add
		cache.Get("boo")
		println("4")
	}()
	wg.Wait() // эта функция ждет пока все рутины завершат работу
	//тут заканчивается код который уже работат с использование Mutex
	*/

	//тут код который уже работает с шардированием
	cache := NewInMemoryCache(5)
	cache.Set("foo", "bar")
	cache.Set("foom", "bard")
	cache.Set("boo", "bah")
	cache.Set("boom", "baha")

}
