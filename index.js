//ADIÇÃO DO MAPA
async function initMap() {
    
    // Função para criar um infoWindow com botões dinâmicos
    function createInfoWindow(content, marker, map) {
        const infoWindow = new google.maps.InfoWindow({
            content: content
        });
        
        // Adiciona um evento de clique para abrir o infoWindow
        marker.addListener("click", () => {
            infoWindow.open(map, marker);
        });
        
        // Configura os eventos dos botões dentro do infoWindow
        infoWindow.addListener('domready', function() {
            // Adiciona um evento de clique para "Traçar Rota"
            document.getElementById('btnRota').addEventListener('click', function(event) {
                event.preventDefault(); // Evita a ação padrão do navegador
                // Função para traçar rota
                calculateAndDisplayRoute(marker.getPosition(), map);
                // Fecha o InfoWindow após clicar
                infoWindow.close();
            });
            
            // Adiciona um evento de clique para "Abrir no Google Maps"
            document.getElementById('btnGoogleMaps').addEventListener('click', function(event) {
                event.preventDefault(); // Evita a ação padrão do navegador
                // Abre o endereço no Google Maps
                const googleMapsUrl = `https://www.google.com/maps/search/?api=1&query=${encodeURIComponent(marker.getPosition().toUrlValue())}`;
                window.open(googleMapsUrl, '_blank');
                // Fecha o InfoWindow após clicar
                infoWindow.close();
            });
        });
        
        return infoWindow;
    }
    
    if (document.getElementById('map')) {
        var map = new google.maps.Map(document.getElementById('map'), {
            center: {lat: -16.6917438, lng: -49.2649191},
            zoom: 13.21
        });
        zoomMapLoc(map);

        directionsService = new google.maps.DirectionsService();
        directionsRenderer = new google.maps.DirectionsRenderer({
            suppressMarkers : true
        });
        directionsRenderer.setMap(map);

        try {
            const response = await fetch('/mapa?json=true');
            const data = await response.json();

            for (const endereco of data.Enderecos) {
                const geocoder = new google.maps.Geocoder();
                const address = `${endereco.Rua}, ${endereco.Numero}, ${endereco.Bairro} ${endereco.Cidade}, ${endereco.CEP}`;
                const nome = `${endereco.Nome}`;

                const results = await geocodeAddress(geocoder, address);

                if (results && results.length > 0) {
                    const icone = {
                        url: "assets/pointer_icon.png",
                        scaledSize: new google.maps.Size(30, 30),
                    };

                    const marker = new google.maps.Marker({
                        map: map,
                        position: results[0].geometry.location,
                        title: nome + "\n" + address,
                        animation: google.maps.Animation.DROP,
                        icon: icone
                    });

                    // Cria o conteúdo do infoWindow dinamicamente
                    const infoWindowContent = `
                        <div>
                            <h2>Opções do Marcador</h2>
                            <button id="btnRota">Traçar Rota</button>
                            <button id="btnGoogleMaps">Abrir no Google Maps</button>
                        </div>
                    `;

                    // Cria e associa o infoWindow ao marcador
                    const infoWindow = createInfoWindow(infoWindowContent, marker, map);
                } else {
                    console.error('Geocode error: ' + status);
                }
            }
        } catch (error) {
            console.error('Erro ao buscar endereços:', error);
        }
    }

    if (document.getElementById('map_paciente')) {
        var geocoder = new google.maps.Geocoder();
        var nome = document.getElementById('nome').innerText
        var endereco = document.getElementById('cidade').innerText + ", " + document.getElementById('cep').innerText + ", " + document.getElementById('bairro').innerText + ", " + document.getElementById('rua').innerText + ", " + document.getElementById('numero').innerText;

        directionsService = new google.maps.DirectionsService();
        directionsRenderer = new google.maps.DirectionsRenderer({
            suppressMarkers: true
        });

        geocoder.geocode({ address: endereco }, function(results, status) {
            if (status === google.maps.GeocoderStatus.OK) {
                var latitude = results[0].geometry.location.lat();
                var longitude = results[0].geometry.location.lng();

                var map_paciente = new google.maps.Map(document.getElementById('map_paciente'), {
                    center: {lat: latitude, lng: longitude},
                    zoom: 13.21
                });

                directionsRenderer.setMap(map_paciente);

                const icone = {
                    url: "assets/icon_house.png",
                    scaledSize: new google.maps.Size(25, 30),
                };

                const paciente = new google.maps.Marker({
                    map: map_paciente,
                    position: results[0].geometry.location,
                    title: nome,
                    animation: google.maps.Animation.DROP,
                    icon: icone
                });

                // Cria o conteúdo do infoWindow para o marcador do paciente
                const infoWindowContent = `
                    <div>
                        <h1 style="position:relative;margin-bottom:40px" >Opções do Marcador</h1>
                        <button id="btnRota">Traçar Rota</button>
                        <button id="btnGoogleMaps">Abrir no Google Maps</button>
                    </div>
                `;

                // Cria e associa o infoWindow ao marcador do paciente
                const infoWindowPaciente = createInfoWindow(infoWindowContent, paciente, map_paciente);
            } else {
                console.error("Geocoding falhou:", status);
            }
        });        
    }
}

async function geocodeAddress(geocoder, address) {
    return new Promise((resolve, reject) => {
        geocoder.geocode({ address }, (results, status) => {
            if (status === 'OK') {
                resolve(results);
            } else {
                reject('Geocode error: ' + status);
            }
        });
    });
}

function zoomMapLoc(map) {
        document.getElementById('noroeste').addEventListener('click', function(event) {
            map.panTo({lat: -16.624041, lng: -49.335088});
        });
        document.getElementById('norte').addEventListener('click', function(event) {
            map.panTo({lat: -16.621392, lng: -49.282103});
        });
        document.getElementById('centro').addEventListener('click', function(event) {
            map.panTo({lat: -16.654146, lng: -49.263644});
        });
        document.getElementById('oeste').addEventListener('click', function(event) {
            map.panTo({lat: -16.690508, lng: -49.373100});
        });
        document.getElementById('leste').addEventListener('click', function(event) {
            map.panTo({lat: -16.695559, lng: -49.182288});
        });
        document.getElementById('sul').addEventListener('click', function(event) {
            map.panTo({lat: -16.728938, lng: -49.277995});
        });
        document.getElementById('sudeste').addEventListener('click', function(event) {
            map.panTo({lat: -16.739157, lng: -49.203613});
        });
        document.getElementById('sudoeste').addEventListener('click', function(event) {
            map.panTo({lat: -16.751362, lng: -49.359509});
        });
}

function calculateAndDisplayRoute(destination, mapa) {
    if (navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(position => {
        const origin = {
          lat: position.coords.latitude,
          lng: position.coords.longitude,
        };

        const icon = {
            url: "assets/acs_icon.png",
            scaledSize: new google.maps.Size(25, 25),
        };

        const acs = new google.maps.Marker({
            map: mapa,
            position: origin,
            animation: google.maps.Animation.DROP,
            icon: icon,
        });

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