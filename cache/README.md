# Custom Cache

## Objective
Design a custom Cache based on TTL with proper expiry of the keys.

## Usecase considered
1. There is a random number generator, which keeps generating numbers between 1 to 100 every second.
2. The random numbers generated are stored in the custom cahce and it has a TTL (time to live) as 5 seconds.
3. There is a random number consumer, which every 30 seconds becomes active and then tries to query some numbers (here it is 5). In a particular iteration, for each of the queried number, if the number is found then print it else say it is a miss.


## Assumptions
1. The items in the cache will be deleted on expiry -> saves memory
2. Expiry check will be done at 2 level:
   1. Active check: Which is a kind of regular check on the cache to clean the expired keys.
   2. Passive check: Which checks if the key is present then whether it is expired or not.
3. With both active and passive check we can achieve optimised cleanup.
4. We can not run active check for long as it will hold the lock, so we will try to keep the active check runtime as less.
