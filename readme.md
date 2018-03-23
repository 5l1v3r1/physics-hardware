<img src="https://s9.postimg.org/9qm3kmdr3/logo4.png" width="300">

Physics platform is a tool for hardware systems (e.g: **raspberryPi 3B**).
It retrieves data passing through the network and sends it to a control panel.
It works the same way as a botnet by receiving remote commands.
(you can imagine that as a black box)

```
  go get github.com/graniet/physics-hardware
```

## command & control

You can check [repository](https://github.com/graniet/physics-command) of command & control

## requirement

You can download a released [binary](https://github.com/graniet/physics-hardware/releases) or create your custom with source code, please make sure you have correctly configured golang environnement with **Go >= 1.8**

  + apt-get install libpcap0.8-dev
  + go get github.com/bettercap/bettercap
  + go get github.com/distatus/battery
  
## wiki

You can check documentation of project [here](https://github.com/graniet/physics-hardware/wiki)
