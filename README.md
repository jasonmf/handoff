# Handoff

Handoff is a tool to easily exchange secret data with no prior key exchange.

Alice wants to receive a secret Bob's cell number from him.

* Alice generates a key with some kind of helpful reference, like `bob-cell`. The tool provides a key: `UvgEzK+TxDKI1pquhg3d7HkI57gw7hyEbYSLQ/62kHhBQkNERUYtMTIz` and saves the key pair it has generated to disk.
* She sends that value to Bob via chat.
* Bob provides the key he received from Alice to the handoff tool. The tool notes that this key is in reference to `bob-cell`. That's what he's expecting to provide so he enters his cell number into the tool. It prints out some encrypted text:

```
-----BEGIN HANDOFF MESSAGE-----
Reference: bob-cell

SXWjozRsbQ14DnlzxcGf7W2geoVVGOjnnHzXWIvWFzbZTi35446bacaNbRF7Rp8d
X8tlIC8JaMAD/1tEdw==
-----END HANDOFF MESSAGE-----
```

* He sends that message back to Alice via chat.
* Alice pastes that value into handoff. Based on the reference ID the handoff tool reads the correct key pair from the disk and recovers Bob's cell number.

Keys in handoff are intended to be single-use. This is not for security purposes but instead to clarify what a secret value is in references.

## Strengths

* Strong encryption ([NACL Box](https://nacl.cr.yp.to/box.html))
* Small keys
* Easy to copy-paste values over chat
* Pure go
* No dependencies; the binary is all you need

## Shortcomings

* Intended for use with small values (4kiB or less)
* No protection against Man in the Middle
* No concept of key trust (e.g. PGP)

## Usage

### CLI

```
alice@blarg:~/p/handoff/cmd$ ./handoff -generate bob-cell
send to user: UvgEzK+TxDKI1pquhg3d7HkI57gw7hyEbYSLQ/62kHhBQkNERUYtMTIz

bob@foo:~/p/handoff/cmd$ ./handoff -encrypt UvgEzK+TxDKI1pquhg3d7HkI57gw7hyEbYSLQ/62kHhBQkNERUYtMTIz
reading for bob-cell from STDIN
555-123-4567
send in response:
-----BEGIN HANDOFF MESSAGE-----
Reference: bob-cell

SXWjozRsbQ14DnlzxcGf7W2geoVVGOjnnHzXWIvWFzbZTi35446bacaNbRF7Rp8d
X8tlIC8JaMAD/1tEdw==
-----END HANDOFF MESSAGE-----

alice@blarg:~/p/handoff/cmd$ ./handoff -decrypt cell-number.txt
reading message from STDIN, writing to cell-number.txt
-----BEGIN HANDOFF MESSAGE-----
Reference: bob-cell

SXWjozRsbQ14DnlzxcGf7W2geoVVGOjnnHzXWIvWFzbZTi35446bacaNbRF7Rp8d
X8tlIC8JaMAD/1tEdw==
-----END HANDOFF MESSAGE-----

alice@blarg:~/p/handoff/cmd$ cat cell-number.txt 
555-123-4567
```

