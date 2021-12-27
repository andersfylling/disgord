This examples will delete its own messages after N seconds, as configured. This examples also showcases middlewares, as these might be of interest to keep your handlers as clean as possible. While they do add some extra code, they become easy to re-use and makes handler registration more readable.

Here global variables are used to store states, which I do currently consider a drawback of middlewares. But since it's easy to make
