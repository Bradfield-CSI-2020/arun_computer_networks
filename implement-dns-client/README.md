# DNS Client

## TODO

1. Decode compressed Name
2. Use random lookup address
3. Multiple questions?
4. Handle response errors 

## Usage

```
go run . google.com
```

Example Output

```
-------Query Headers -------:
Reply ID: 46264
QR Flag: 0
OPCode Flag: 0
AA Flag: 0
TC Flag: 0
RD Flag: 0
RA Flag: 0
RCode Flag: 0
No. of Questions: 1
No. of Answers: 0
No. of Name Servers: 0
No. of Authoritative Records: 0

-------Query Question -------:
QName: %!d(string=yahoo.com)
QClass: 0
QType: 0
size of query:  27
size of response:  123

-------Reply Headers -------:
Reply ID: 46264
QR Flag: 1
OPCode Flag: 0
AA Flag: 0
TC Flag: 0
RD Flag: 0
RA Flag: 1
RCode Flag: 0
No. of Questions: 1
No. of Answers: 6
No. of Name Servers: 0
No. of Authoritative Records: 0

-------Reply Ansswer -------:
Name: [192 12]
Type: [0 1]
Class: [0 1]
TTL: 1482
RD Length: 4
IP Address:  98.137.11.163
```
