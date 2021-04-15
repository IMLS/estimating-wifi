# Installation and configuration

This is the exciting step. It is also kinda slow.

Note that once you follow these steps, you will not be able to access your Raspberry Pi again. For security reasons, we lock all users out of the device. You can "reclaim" your Pi by removing the microSD card and re-imaging it (starting over, basically).

# Installation

In the upper-left-hand corner of the Pi's display, in the title bar, is an icon for a little screen.

{% asset "pi-command-prompt.png" alt="location of pi command prompt icon" %}

Click this. It will open a terminal. Now, paste in the following line.

```
bash <$(curl -s https://raw.githubusercontent.com/jadudm/imls-client-pi-playbook/main/bootstrap.sh)
```

FIXME: more docs here.

# Configuration

FIXME: more here

[The token translator](token.md)

