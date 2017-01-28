package single

import (
	"errors"
	"testing"
	"time"

	. "github.com/jochasinga/grx/bases"
	"github.com/jochasinga/grx/handlers"
	"github.com/jochasinga/grx/observer"
	"github.com/jochasinga/grx/subscription"

	"github.com/stretchr/testify/assert"
)

type (
	Number int
	Text   string
)

func (num Number) Emit() (Item, error) {
	return Item(num), nil
}

func (tx Text) Emit() (Item, error) {
	return (Item)(nil), errors.New("text error")
}

func TestSingleImplementIterator(t *testing.T) {
	assert.Implements(t, (*Iterator)(nil), DefaultSingle)
}

func TestSingleImplementStream(t *testing.T) {
	assert.Implements(t, (*Stream)(nil), DefaultSingle)
}

func TestCreateSingleWithConstructor(t *testing.T) {
	assert := assert.New(t)

	s := New(Number(1))

	emitter, err := s.Next()
	assert.Nil(err)
	assert.NotNil(emitter)
	assert.Implements((*Emitter)(nil), emitter)
	assert.EqualValues(1, emitter)

	emitter, err = s.Next()
	assert.Nil(emitter)
	assert.NotNil(err)
	assert.Panics(func() {
		assert.Implements((*Emitter)(nil), emitter)
	})
	assert.EqualValues(nil, emitter)
}

func TestSubscribingToObserver(t *testing.T) {
	assert := assert.New(t)
	num := 2
	errorMessage := ""
	ob := &observer.Observer{
		NextHandler: handlers.NextFunc(func(item Item) {
			num += int(item.(Number))
		}),
		ErrHandler: handlers.ErrFunc(func(err error) {
			errorMessage = err.Error()
		}),
	}

	s1 := New(Number(1))
	sub, err := s1.Subscribe(ob)
	<-time.After(10 * time.Millisecond)
	if s, ok := sub.(*subscription.Subscription); ok {
		assert.WithinDuration(s.SubscribeAt(), time.Now(), 20*time.Millisecond)
	}
	assert.Nil(err)
	assert.Implements((*Subscriptor)(nil), sub)
	assert.Equal(3, num)

	s2 := New(Text("Hello"))
	sub, err = s2.Subscribe(ob)
	<-time.After(10 * time.Millisecond)
	assert.Nil(err)
	assert.Implements((*Subscriptor)(nil), sub)
	assert.Equal("text error", errorMessage)
}

func TestSubscribingToHandlers(t *testing.T) {
	assert := assert.New(t)
	num := 2
	errorMessage := ""

	nextf := handlers.NextFunc(func(item Item) {
		num += int(item.(Number))
	})
	errf := handlers.ErrFunc(func(err error) {
		errorMessage = err.Error()
	})

	s1 := New(Number(1))
	sub, err := s1.Subscribe(nextf)
	<-time.After(10 * time.Millisecond)
	if s, ok := sub.(*subscription.Subscription); ok {
		assert.WithinDuration(s.SubscribeAt(), time.Now(), 20*time.Millisecond)
	}
	assert.Nil(err)
	assert.Implements((*Subscriptor)(nil), sub)
	assert.Equal(3, num)

	s2 := New(Text("Hello"))
	sub, err = s2.Subscribe(errf)
	<-time.After(10 * time.Millisecond)
	assert.Nil(err)
	assert.Implements((*Subscriptor)(nil), sub)
	assert.Equal("text error", errorMessage)
}

func TestSubscribeAndUnsubscribeReturnSameSubscriptor(t *testing.T) {
	assert := assert.New(t)
	s1 := DefaultSingle

	sampleHandler := handlers.NextFunc(func(item Item) {
		return
	})

	sub1, _ := s1.Subscribe(sampleHandler)
	<-time.After(10 * time.Millisecond)
	sub2 := s1.Unsubscribe()

	assert.Equal(sub1.SubscribeAt(), sub2.SubscribeAt())
	assert.Equal(sub1.UnsubscribeAt(), sub2.UnsubscribeAt())
}