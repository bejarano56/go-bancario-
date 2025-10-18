function cargarCuentas() {
  console.log('Iniciando cargarCuentas...'); // Debug
  fetch('/cuentas')
    .then(r => r.json())
    .then(data => {
      console.log('Datos recibidos:', data); // Debug
      const tbody = document.querySelector('#tabla-cuentas tbody');
      tbody.innerHTML = '';
      if (Array.isArray(data)) {
        data.forEach(c => {
          console.log('Procesando cuenta:', c); // Debug
          // Manejar saldo de forma segura
          const saldo = c.SaldoActual || 0;
          const saldoCOP = new Intl.NumberFormat('es-CO', { style: 'currency', currency: 'COP', maximumFractionDigits: 0 }).format(saldo);
          // Usar el nombre del estado si está disponible, sino mostrar el ID
          const estadoText = c.EstadoNombre || `Estado ${c.IDEstado}`;
          // Usar el nombre del tipo si está disponible, sino mostrar el ID
          const tipoText = c.TipoNombre || `Tipo ${c.IDTipoCuenta}`;
          console.log('Valores finales:', { saldoCOP, estadoText, tipoText }); // Debug
          tbody.innerHTML += `<tr><td>${c.IDCuenta}</td><td>${c.NumeroCuenta}</td><td>${c.IDCliente}</td><td>${tipoText}</td><td>${saldoCOP}</td><td>${estadoText}</td></tr>`;
        });
      }
    })
    .catch(error => {
      console.error('Error cargando cuentas:', error);
    });
}

document.addEventListener('DOMContentLoaded', () => {
  console.log('DOM cargado, iniciando...'); // Debug
  cargarCuentas();
});


