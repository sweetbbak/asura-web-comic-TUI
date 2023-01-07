# asura scans web-comic TUI

Search the front page of the popular manhwa / webcomic site, asura scans and then view the chapter right in your terminal (kitty only)

![gif demonstration](https://github.com/sweetbbak/asura-web-comic-TUI/blob/main/asura_testes.gif)

# Requirements: 
(beautiful soup 4 is the main one, the other dependencies can be changed to "read" or fzf and you can use the native kitty +kitten icat function)

charmbracelet/gum

```
sudo pacman -S gum && pip install bs4 && pip install pixcat
```

```
 @@@@@@    @@@@@@   @@@  @@@  @@@@@@@    @@@@@@   
@@@@@@@@  @@@@@@@   @@@  @@@  @@@@@@@@  @@@@@@@@  
@@!  @@@  !@@       @@!  @@@  @@!  @@@  @@!  @@@  
!@!  @!@  !@!       !@!  @!@  !@!  @!@  !@!  @!@  
@!@!@!@!  !!@@!!    @!@  !@!  @!@!!@!   @!@!@!@!  
!!!@!!!!   !!@!!!   !@!  !!!  !!@!@!    !!!@!!!!  
!!:  !!!       !:!  !!:  !!!  !!: :!!   !!:  !!!  
:!:  !:!      !:!   :!:  !:!  :!:  !:!  :!:  !:!  
::   :::  :::: ::   ::::: ::  ::   :::  ::   :::  
 :   : :  :: : :     : :  :    :   : :   :   : :  


```

feel free to send a PR!

# To do:
- images currently load at the bottom requiring you to scroll to the top & then back down (#1 priority)
- fix download option
- add the ability to remember your place and track read comics
- add favorites list
- make the menu more interactive and intuitive
- add caching
- support external readers (nsxiv, sxiv, mangreader)