# Golang Presidential-Action Tracker

This bot will track the latest presidential actions listed on [whitehouse.gov](https://www.whitehouse.gov/) and alert users of new Exective Orders.

# Installation

To run the bot, run the following commands:

```
git clone https://github.com/NoahSoto/executive-order-tracker
cd ./executive-order-tracker
pip install -r requirements.txt
go run main.go &
```

Or, if you prefer to build an executable that will run without the overhead, run:

```
git clone https://github.com/NoahSoto/executive-order-tracker
cd ./executive-order-tracker
pip install -r requirements.txt
go build -o executive-order-tracker-binary main.go
./executive-order-tracker-binary
```

Note that the executable option will consume more storage but require less headroom, as well as require a rebuild every time the main source code is updated.

# Discord Usage

In Discord, to search for executive orders, run the command:

```
!orders
```

You can view previous orders by running: 

```
!ls
```

Finally, you can retrieve the title, URL, and a brief NLTK-based synopsis with the command:

```
!view [1-10]
```

If you have questions you can always run:

```
!help
```

