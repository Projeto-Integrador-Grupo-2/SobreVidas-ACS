let map;

async function initMap() {
  const { Map } = await google.maps.importLibrary("maps");

  map = new Map(document.getElementById("map"), {
    center: { lat: -16.6917438, lng: -49.2649191},
    zoom: 13.21,
  });
}

initMap();

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


/* Editar paciente*/
function editPatient(event, id) {
    event.preventDefault();
    window.location.href = `/cadastro.html?id=${id}`;
}

document.addEventListener('DOMContentLoaded', function() {
    const urlParams = new URLSearchParams(window.location.search);
    const patientId = urlParams.get('id');

    if (patientId) {
        console.log(`Buscando dados do paciente com ID: ${patientId}`);
        fetch(`/getPaciente?id=${patientId}`)
            .then(response => {
                if (!response.ok) {
                    throw new Error('Erro ao buscar dados do paciente');
                }
                return response.json();
            })
            .then(data => {
                console.log('Dados do paciente:', data);
                
                
                // Verifique se cada campo está presente e defina um valor padrão se não estiver
                document.getElementById('patientId').value = data.Id || '';
                document.getElementById('data_cadastro').value = formatarData(data.Data_cad);
                document.getElementById('nome').value = data.nome || '';
                document.getElementById('nome_mae').value = data.Nome_mae || '';
                document.getElementById('cpf').value = data.cpf || '';
                document.getElementById('sexo').value = data.sexo || '';
                document.getElementById('email').value = data.email || '';
                document.getElementById('telefone').value = data.celular || '';
                document.getElementById('data_nascimento').value = formatarData(data.data_nasc);
                document.getElementById('cidade').value = data.Cidade || '';
                document.getElementById('cep').value = data.CEP || '';
                document.getElementById('logradouro').value = data.Rua || '';
                document.getElementById('numero').value = data.Num_casa || '';
            })
            .catch(error => {
                console.error('Erro ao carregar dados do paciente:', error);
            });
    }
});

function formatarData(data) {
    var partes = data.split('/');
    return partes[0] + '-' + partes[1] + '-' + partes[2];
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