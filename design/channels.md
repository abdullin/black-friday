# Channels


This is a fun demo concept.

We create a sales channel to serve things. Quantity of a product available for a sales channel
can be different from on-hand quantity.

It will depend:
- how much stock was reserved for the channel (is it a campaign), also at which price
- how much stock can be back-ordered.
- if the stock can be moved between the warehouses


Sales are placed via channels. They usually come with SKUs, quantities and prices.

When sale comes, we reserve onHand stock for it.

Each channel is a stateful executable application that can deploy stocks (counters and allocations)




