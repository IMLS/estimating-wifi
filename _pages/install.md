---
title: Install and Config
layout: page
sidenav: false
---
{% assign counter = 0 %}

{{page.title}}

Here we go! You have an RPi, it is set up, and it can talk to the internet! This is very, *very* exciting.

## Step {% assign counter = counter | plus: 1  %}{{counter}}: Navigate to this page on the RPi

We have created a short URL for this page. This will make it easier to get here.

FIXME

In the web browser on the RPi, type that URL in. You should find yourself *right here*. 

## Step {% assign counter = counter | plus: 1  %}{{counter}}: Open the terminal

In the upper-left corner of the screen is a little black box with a blue bar at the top. This icon is intended to represent a command terminal. It should open up a window that looks like this:

{% asset "pi-command-prompt.png" alt="location of pi command prompt icon" %}

## Step {% assign counter = counter | plus: 1  %}{{counter}}: Copy the line below

```
bash <(curl -s https://raw.githubusercontent.com/jadudm/imls-client-pi-playbook/main/bootstrap.sh)>
```

Highlight that line, right-click, and say **Copy**. 

## Step {% assign counter = counter | plus: 1  %}{{counter}}: Paste the line into the terminal

Over in the terminal, right-click, and select **Paste**. The command should now be in the terminal window.

Press Enter. Or, Return. Or whatever that key is called on your keyboard.

## Step {% assign counter = counter | plus: 1  %}{{counter}}: Enter your FCFS_Seq_Id

You are setting up your device so that it will live in a particular library. Or, perhaps, a bookmobile. Either way, there should be an FSFS_Seq_Id for your location.

You can [look this up on the IMLS website](https://www.imls.gov/search-compare/). For example, the Berea branch in Madison County, KY is **KY0069-003**.

## Step {% assign counter = counter | plus: 1  %}{{counter}}: Enter a tag

This is a free-form text entry that will accept up to 255 characters. It *will* become part of the dataset, so please do not put any personally identifiable information into this field. Email addresses are considered PII.

We intend this field either as a "hardware tag" (so you might enter "RPi 001"), or it could be a reminder where you put it. In the latter case, you might enter "reference desk," or "networking closet," or (if it is in the bookmobile)... er, "glovebox?" Point being, this is mostly for you. Again, this will become part of the dataset, so something simple but descriptive is best.

## Step {% assign counter = counter | plus: 1  %}{{counter}}: Enter your api.data.gov key

This is a tedious part.




## Step {% assign counter = counter | plus: 1  %}{{counter}}: Installation and configuration

This is the exciting step. It is also kinda slow.

Note that once you follow these steps, you will not be able to access your Raspberry Pi again. For security reasons, we lock all users out of the device. You can "reclaim" your Pi by removing the microSD card and re-imaging it (starting over, basically).

# Installation

In the upper-left-hand corner of the Pi's display, in the title bar, is an icon for a little screen.

{% asset "pi-command-prompt.png" alt="location of pi command prompt icon" %}

Click this. It will open a terminal. Now, paste in the following line.

```
bash <(curl -s https://raw.githubusercontent.com/jadudm/imls-client-pi-playbook/main/bootstrap.sh)
```

FIXME: more docs here.

# Configuration

FIXME: more here

[The token translator](token.md)

