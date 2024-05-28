const menu = document.getElementById('opcoes_icon')
const menucontent = document.getElementById('opcoes-content')
let menuaberto = false;

menu.addEventListener('click', () => {
    if (!menuaberto) {
        menucontent.style.display = 'block'
        menuaberto = true;
    }
    else {
        menucontent.style.display = 'none';
        menuaberto = false;
    }
    
})
document.addEventListener("click", function(event) {
    if (!menucontent.contains(event.target) && event.target !== menu) {
      menucontent.style.display = 'none';
      menuaberto = false;
    }
})