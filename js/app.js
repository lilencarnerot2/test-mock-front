const API_URL = "http://localhost:8092/alojamientos";

let map;
let globalData = [];

// LOGIN
document.getElementById("btnLogin").addEventListener("click", () => {
  const user = document.getElementById("usuario").value;
  const pass = document.getElementById("password").value;

  fetch(API_URL, {
    headers: {
      "Authorization": "Basic " + btoa(user + ":" + pass)
    }
  })
    .then(res => res.json())
    .then(data => {
      globalData = data.data;

      render(globalData);
      renderMap(globalData);
    })
    .catch(err => console.log(err));
});

// OJITO PASSWORD
document.getElementById("togglePass").addEventListener("click", () => {
  const pass = document.getElementById("password");
  pass.type = pass.type === "password" ? "text" : "password";
});

// RESET
document.getElementById("btnReset").addEventListener("click", () => {
  document.getElementById("app").innerHTML = "";
  if (map) {
    map.remove();
    map = null;
  }
  globalData = [];
});

// SEARCH
document.getElementById("search").addEventListener("input", (e) => {
  const value = e.target.value.toLowerCase();

  const filtered = globalData.filter(item =>
    item.nombre.toLowerCase().includes(value)
  );

  render(filtered);
  renderMap(filtered);
});

// FILTER
document.getElementById("filter").addEventListener("change", (e) => {
  const value = e.target.value;

  if (value === "all") {
    render(globalData);
    renderMap(globalData);
    return;
  }

  const filtered = globalData.filter(i => i.tipo === value);

  render(filtered);
  renderMap(filtered);
});

// RENDER CARDS
function render(data) {
  document.getElementById("app").innerHTML = `
    <h2>Total: ${data.length}</h2>

    <div class="cards">
      ${data.map(i => `
        <div class="card">
          <h3>${i.nombre}</h3>
          <p>${i.domicilio}</p>
          <small>${i.tipo}</small>
        </div>
      `).join("")}
    </div>
  `;
}

// MAPA
function renderMap(data) {
  if (!map) {
    map = L.map('map').setView([-34.66, -58.72], 12);

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: 'OpenStreetMap'
    }).addTo(map);
  }

  map.eachLayer(layer => {
    if (layer instanceof L.Marker) {
      map.removeLayer(layer);
    }
  });

  data.forEach(item => {
    if (item.lat && item.long) {
      L.marker([item.lat, item.long])
        .addTo(map)
        .bindPopup(item.nombre);
    }
  });
}