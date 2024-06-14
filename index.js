document.addEventListener('DOMContentLoaded', function() {
    var cpfInput = document.getElementById('cpf');
    var telefoneInput = document.getElementById('telefone');
    var cepInput = document.getElementById('cep');

    VMasker(cpfInput).maskPattern('999.999.999-99');
    VMasker(telefoneInput).maskPattern('(99) 99999-9999');
    VMasker(cepInput).maskPattern('99999-999');
});


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

function toggleDropdown(event) {
    event.stopPropagation();
    const dropdownContent = event.currentTarget.nextElementSibling;
    dropdownContent.classList.toggle("show");
}

window.onclick = function(event) {
    if (!event.target.matches('.dropdown-button')) {
        const dropdowns = document.getElementsByClassName("dropdown-content");
        for (let i = 0; i < dropdowns.length; i++) {
            const openDropdown = dropdowns[i];
            if (openDropdown.classList.contains('show')) {
                openDropdown.classList.remove('show');
            }
        }
    }
}

let deletePatientId = null;

function toggleDropdown(event) {
    event.stopPropagation();
    const dropdownContent = event.currentTarget.nextElementSibling;
    dropdownContent.classList.toggle("show");
}

window.onclick = function(event) {
    if (!event.target.matches('.dropdown-button')) {
        const dropdowns = document.getElementsByClassName("dropdown-content");
        for (let i = 0; i < dropdowns.length; i++) {
            const openDropdown = dropdowns[i];
            if (openDropdown.classList.contains('show')) {
                openDropdown.classList.remove('show');
            }
        }
    }
}

function confirmDelete(event, id) {
    event.preventDefault();
    console.log(`Paciente a ser excluído com ID: ${id}`); // Log do ID do paciente
    deletePatientId = id;
    document.getElementById("deleteModal").style.display = "block";
}

function closeModal() {
    document.getElementById("deleteModal").style.display = "none";
    deletePatientId = null;
}

function deletePatient() {
    if (deletePatientId === null) return;

    fetch(`/deletePaciente`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ id: deletePatientId })
    }).then(response => {
        if (response.ok) {
            location.reload();
        } else {
            alert("Erro ao excluir o paciente.");
        }
    }).catch(error => {
        console.error("Erro:", error);
        alert("Erro ao excluir o paciente.");
    }).finally(() => {
        closeModal();
    });
}

const slides = document.querySelector('.slides');
const slideCount = document.querySelectorAll('.slide').length;
const navDots = document.querySelectorAll('.nav-dot');
let currentIndex = 0;
let slideInterval;

function nextSlide() {
    currentIndex = (currentIndex + 1) % slideCount;
    updateSlidePosition();
}

function updateSlidePosition() {
    const offset = -currentIndex * 100;
    slides.style.transform = `translateX(${offset}%)`;
    updateNavDots();
}

function updateNavDots() {
    navDots.forEach(dot => dot.classList.remove('active'));
    navDots[currentIndex].classList.add('active');
}

function goToSlide(index) {
    currentIndex = index;
    updateSlidePosition();
    resetInterval();
}

function resetInterval() {
    clearInterval(slideInterval);
    slideInterval = setInterval(nextSlide, 15000);
}

navDots.forEach(dot => {
    dot.addEventListener('click', () => {
        const index = parseInt(dot.getAttribute('data-index'));
        goToSlide(index);
    });
});

slideInterval = setInterval(nextSlide, 10000);
updateNavDots();