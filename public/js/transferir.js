document.addEventListener('DOMContentLoaded', () => {
  const form = document.getElementById('form-transferir');
  if (form) {
    form.addEventListener('submit', (e) => {
      e.preventDefault();
      const id_origen = document.getElementById('id_origen').value;
      const id_destino = document.getElementById('id_destino').value;
      const id_tipo_transaccion = document.getElementById('id_tipo_transaccion').value;
      const monto = document.getElementById('monto').value;
      let referencia = document.getElementById('referencia').value;
      const descripcion = document.getElementById('descripcion').value;
      fetch('/cuentas', {
        method: 'PUT',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: `id_origen=${encodeURIComponent(id_origen)}&id_destino=${encodeURIComponent(id_destino)}&id_tipo_transaccion=${encodeURIComponent(id_tipo_transaccion)}&monto=${encodeURIComponent(monto)}&referencia=${encodeURIComponent(referencia)}&descripcion=${encodeURIComponent(descripcion)}`
      }).then(resp => {
        return resp.text().then(text => {
          const alert = document.createElement('div');
          alert.className = `alert ${resp.ok ? 'alert-success' : 'alert-danger'} mt-2`;
          alert.textContent = resp.ok ? 'Transferencia realizada.' : 'Error en la transferencia: ' + text;
          form.appendChild(alert);
          setTimeout(() => alert.remove(), 5000);
          if (resp.ok) form.reset();
        });
      });
    });
  }
});


