//ADIÇÃO DO MAPA
let map;

async function initMap() {
    var map = new google.maps.Map(document.getElementById('map'), {
        center: {lat: -16.6917438, lng: -49.2649191},
        zoom: 13.21
    });

    directionsService = new google.maps.DirectionsService();
    directionsRenderer = new google.maps.DirectionsRenderer();
    directionsRenderer.setMap(map);

    fetch('/mapa?json=true')
        .then(response => response.json())
        .then(data => {
            data.Enderecos.forEach(endereco => {
                // Geocode para obter coordenadas
                var geocoder = new google.maps.Geocoder();
                var address = `${endereco.Rua}, ${endereco.Numero}, ${endereco.Bairro} ${endereco.Cidade}, ${endereco.CEP}`;
                var nome = `${endereco.Nome}`
                        
                geocoder.geocode({ 'address': address }, function(results, status) {
                    if (status === 'OK') {
                        const icone = {
                            url: "assets/pointer_icon.png",
                            scaledSize: new google.maps.Size(22, 22),
                        };

                        const marker = new google.maps.Marker({
                            map: map,
                            position: results[0].geometry.location,
                            title: nome + "\n" + address,
                            animation: google.maps.Animation.DROP,
                            icon: icone
                        });

                        marker.addListener("click", () => {
                            calculateAndDisplayRoute(marker.getPosition());
                        });

                } else {
                    console.error('Geocode error: ' + status);
                }
            });
         });
    })
    .catch(error => console.error('Erro ao buscar endereços:', error));
}

function calculateAndDisplayRoute(destination) {
    if (navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(position => {
        const origin = {
          lat: position.coords.latitude,
          lng: position.coords.longitude,
        };

        directionsService.route(
          {
            origin: origin,
            destination: destination,
            travelMode: google.maps.TravelMode.DRIVING,
          },
          (response, status) => {
            if (status === "OK") {
              directionsRenderer.setDirections(response);
            } else {
              window.alert("Erro ao calcular a rota: " + status);
            }
          }
        );
      });
    } else {
      window.alert("Geolocalização não é suportada pelo navegador.");
    }
  }

//


//MÁSCARAS

document.addEventListener('DOMContentLoaded', function() {
    var cpfInput = document.getElementById('cpf');
    var telefoneInput = document.getElementById('telefone');
    var cepInput = document.getElementById('cep');

    VMasker(cpfInput).maskPattern('999.999.999-99');
    VMasker(telefoneInput).maskPattern('(99) 99999-9999');
    VMasker(cepInput).maskPattern('99999-999');
});
//

//MUDAR TAMANHO DOS GRÁFICOS NO PERFIL
function getElementWidth(selector) {
    const element = document.querySelector(selector);
    if (element) {
        return parseFloat(window.getComputedStyle(element).width);
    }
    return null;
}

const largura1 = getElementWidth('.graf_atendidos');
const largura2 = getElementWidth('.graf_encaminhados');

function mudarTamanho(largura, num, elemento) {
    if (!largura || !num || !elemento) {
        return;
    }
    const numero = parseInt(num.textContent);
    if (isNaN(numero)) {
        return;
    }
    elemento.style.width = largura + numero*0.3 + 'px';
}

if (largura1 !== null) {
    mudarTamanho(largura1, document.getElementById("valor1"), document.getElementById("graf_atend"));
    mudarTamanho(largura2, document.getElementById("valor2"), document.getElementById("graf_encam"));
}

//

//ÍCONE DE OPÇÕES(MENU DROPDOWN)
const menu = document.getElementById('opcoes_icon');
const menucontent = document.getElementById('opcoes-content');


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
//

//FUNÇÃO DELETAR PACIENTE
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
//


//FUNÇÃO EDITAR PACIENTE
function editPatient(event, id) {
    event.preventDefault();
    window.location.href = `/cadastro.html?id=${id}`;
}

document.addEventListener('DOMContentLoaded', function() {
    const urlParams = new URLSearchParams(window.location.search);
    const patientId = urlParams.get('id');

    if (patientId) {
        fetch(`/getPaciente?id=${patientId}`)
            .then(response => {
                if (!response.ok) {
                    throw new Error('Erro ao buscar dados do paciente');
                }
                return response.json();
            })
            .then(data => {                              
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
                document.getElementById('bairro').value = data.Bairro || '';
                document.getElementById('logradouro').value = data.Rua || '';
                document.getElementById('numero').value = data.Num_casa || '';
                document.getElementById('bebe').checked = data.Bebe;
                document.getElementById('fuma').checked = data.Fuma || '';
                document.getElementById('possui_feridas_boca').checked = data.Possui_feridas_boca || '';
            })
            .catch(error => {
                console.error('Erro ao carregar dados do paciente:', error);
            });
    }
});
//

//FORMATAR DADA PARA PADRÃO HTML5
function formatarData(data) {
    var partes = data.split('/');
    return partes[0] + '-' + partes[1] + '-' + partes[2];
}
//


//SLIDESHOW AUTOMÁTICO (CARROSSEL DE IMAGENS)
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
//