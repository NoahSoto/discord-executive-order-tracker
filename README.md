# Golang Presidential-Action Tracker

This bot will track the 10 latest presidential actions listed on [whitehouse.gov](https://www.whitehouse.gov/)

# Install

To install the bot, run the following commands:

```
git clone https://github.com/NoahSoto/executive-order-tracker
cd ./executive-order-tracker
pip install -r requirements.txt
go run main.go &
```

# Discord Usage

In Discord, to initialize the bot and start searching for executive orders, run the command:

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

