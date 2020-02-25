# Distinct Operator

## Overview

Suppress duplicate items emitted by an Observable.

![](http://reactivex.io/documentation/operators/images/distinct.png)

## Example

```go
observable := rxgo.Just(1, 2, 2, 3, 4, 4, 5)().
	Distinct(func(_ context.Context, i interface{}) (interface{}, error) {
		return i, nil
	})
```

Output:

```
1
2
3
4
5
```

## Options

### WithBufferedChannel

[Detail](options.md#withbufferedchannel)

### WithContext

[Detail](options.md#withcontext)

### WithObservationStrategy

[Detail](options.md#withobservationstrategy)

### WithErrorStrategy

[Detail](options.md#witherrorstrategy)

### WithPool

[Detail](options.md#withpool)

### WithCPUPool

[Detail](options.md#withcpupool)

### Serialize

[Detail](options.md#serialize)

### WithPublishStrategy

[Detail](options.md#withpublishstrategy)