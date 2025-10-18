function cargarCuentas() {
  fetch('/cuentas')
    .then(r => r.json())
    .then(data => {
      const tbody = document.querySelector('#tabla-cuentas tbody');
      tbody.innerHTML = '';
      if (Array.isArray(data)) {
        data.forEach(c => {
          const saldoCOP = new Intl.NumberFormat('es-CO', { style: 'currency', currency: 'COP', maximumFractionDigits: 0 }).format(c.SaldoActual);
          tbody.innerHTML += `<tr><td>${c.IDCuenta}</td><td>${c.NumeroCuenta}</td><td>${c.IDCliente}</td><td>${c.IDTipoCuenta}</td><td>${saldoCOP}</td><td>${c.IDEstado}</td></tr>`;
        });
      }
    });
}

document.addEventListener('DOMContentLoaded', () => {
  cargarCuentas();
});


