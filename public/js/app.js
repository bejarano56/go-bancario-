function cargarCuentas() {
  fetch('/cuentas')
    .then(r => r.json())
    .then(data => {
      const tbody = document.querySelector('#tabla-cuentas tbody');
      tbody.innerHTML = '';
      if (Array.isArray(data)) {
        data.forEach(c => {
          const activaText = c.Activa ? 'Sí' : 'No';
          const saldoCOP = new Intl.NumberFormat('es-CO', { style: 'currency', currency: 'COP', maximumFractionDigits: 0 }).format(c.Saldo);
          tbody.innerHTML += `<tr><td>${c.IDCuenta}</td><td>${c.NumeroCuenta}</td><td>${c.IDCliente}</td><td>${c.IDTipoCuenta}</td><td>${saldoCOP}</td><td>${activaText}</td></tr>`;
        });
      }
    });
}

document.addEventListener('DOMContentLoaded', () => {
  // cargar tipos de cuenta
  fetch('/tipos-cuenta').then(r => r.json()).then(tipos => {
    const sel = document.getElementById('id_tipo_cuenta');
    if (sel && Array.isArray(tipos)) {
      sel.innerHTML = '<option value="" disabled selected>Selecciona un tipo</option>';
      tipos.forEach(t => {
        sel.innerHTML += `<option value="${t.IDTipoCuenta}">${t.NombreTipo}</option>`;
      });
    }
  });

  const formCuenta = document.getElementById('form-cuenta');
  if (formCuenta) {
    formCuenta.addEventListener('submit', (e) => {
      e.preventDefault();
      const id_tipo_cuenta = document.getElementById('id_tipo_cuenta').value;
      const nombre_completo = document.getElementById('nombre_completo').value;
      fetch('/cuentas', {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: `nombre_completo=${encodeURIComponent(nombre_completo)}&id_tipo_cuenta=${encodeURIComponent(id_tipo_cuenta)}`
      }).then(resp => {
        if (resp.ok) {
          const alert = document.createElement('div');
          alert.className = 'alert alert-success mt-2';
          alert.textContent = 'Cuenta creada correctamente con saldo inicial de $500.000 COP';
          formCuenta.appendChild(alert);
          setTimeout(() => alert.remove(), 2500);
          cargarCuentas();
          formCuenta.reset();
        }
      });
    });
  }

  const formTransferir = document.getElementById('form-transferir');
  if (formTransferir) {
    formTransferir.addEventListener('submit', (e) => {
      e.preventDefault();
      const id_origen = document.getElementById('id_origen').value;
      const id_destino = document.getElementById('id_destino').value;
      const id_tipo_transaccion = document.getElementById('id_tipo_transaccion').value;
      const monto = document.getElementById('monto').value;
      const referencia = document.getElementById('referencia').value;
      const descripcion = document.getElementById('descripcion').value;
      fetch('/cuentas', {
        method: 'PUT',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: `id_origen=${encodeURIComponent(id_origen)}&id_destino=${encodeURIComponent(id_destino)}&id_tipo_transaccion=${encodeURIComponent(id_tipo_transaccion)}&monto=${encodeURIComponent(monto)}&referencia=${encodeURIComponent(referencia)}&descripcion=${encodeURIComponent(descripcion)}`
      }).then(() => {
        cargarCuentas();
        formTransferir.reset();
      });
    });
  }

  cargarCuentas();
});


