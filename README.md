# Rate Limiter Algorithms

This repository contains implementations of various rate limiting algorithms using Go. These algorithms help manage how often users or clients can access a resource within a given time frame, which is crucial for maintaining the stability and reliability of applications.

## Algorithms Implemented

### 1. **Token Bucket**:

The Token Bucket is metaphorically like a bucket where tokens are continuously added at a predetermined rate. Each token allows for one unit of work or one request to be processed. The bucket has a capacity, which is the maximum number of tokens it can hold.

<p align="center">
    <img style="width:400px" src="https://pub.anw.sh/token-bucket.png"/>
</p>

When a request arrives: The algorithm first checks if there are enough tokens in the bucket to handle the request. Each request typically requires one token, but depending on the system’s design, it might require more.

1. Bucket has tokens: tokens are removed from the bucket, and the request is allowed to proceed.
1. Bucket doesn't have tokens: the request is either denied or queued for later processing, depending on the system’s design.

**Features**:

- Burst Capacity: one key feature of the Token Bucket is its ability to handle bursts of requests up to its capacity. If the bucket is full, it can handle a number of requests equal to its capacity all at once, provided that no more tokens are required until the bucket begins to refill.

### 2. **Leaky Bucket**:

The Leaky Bucket algorithm can be visualized as a bucket where each incoming request adds some amount of water, and there is a hole in the bottom through which water leaks out at a constant rate.

<p align="center">
    <img style="width:400px" src="https://pub.anw.sh/leaky-bucket.png"/>
</p>

When a request arrives: The algorithm checks if there is room in the bucket to add the “water” from the new request without overflowing.

1. Bucket has room: it can accommodate the water from the incoming request, the request is added to the queue (i.e. water is added to the bucket), and the request is processed at the rate at which water leaks out.
2. Bucket has no room: If adding water from the request would cause the bucket to overflow (i.e. exceed its capacity), the request is denied. This is to prevent the bucket from overflowing and ensures that requests are processed at no faster rate than the leak rate allows.

**Features**:

- Steady Throughput: unlike the Token Bucket, the Leaky Bucket enforces a steady rate of requests and doesn’t naturally allow for bursts. If requests come in faster than they can leak out, the excess requests are denied, maintaining a smooth and consistent output rate.

### Other simple algorithms

1. **Fixed Window Counter**: Counts requests in a fixed time window (e.g., per minute) and resets the count at the start of the next window. Can potentially allow bursts of traffic at the edges of the windows.
2. **Sliding Log**: Records timestamps of each incoming request and ensures that the rate of requests doesn't exceed the set limit for any given time window.
3. **Sliding Window Counter**: Improves on the fixed window by smoothing out bursts at the edges of the windows using overlapping windows to calculate the rate.
