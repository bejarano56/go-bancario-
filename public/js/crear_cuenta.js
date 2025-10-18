document.addEventListener('DOMContentLoaded', () => {
  // Cargar tipos de cuenta
  fetch('/tipos-cuenta').then(r => r.json()).then(tipos => {
    const sel = document.getElementById('id_tipo_cuenta');
    if (sel && Array.isArray(tipos)) {
      sel.innerHTML = '<option value="" disabled selected>Selecciona un tipo</option>';
      tipos.forEach(t => {
        const option = document.createElement('option');
        option.value = t.IDTipoCuenta;
        option.textContent = t.NombreTipo;
        sel.appendChild(option);
      });
    }
  }).catch(error => {
    console.error('Error cargando tipos de cuenta:', error);
  });

  const form = document.getElementById('form-cuenta');
  if (form) {
    console.log('Form found, adding submit listener');
    form.addEventListener('submit', (e) => {
      console.log('Form submit triggered');
      e.preventDefault();
      const id_tipo_cuenta = document.getElementById('id_tipo_cuenta').value;
      const tipo_documento = document.getElementById('tipo_documento').value;
      const numero_documento = document.getElementById('numero_documento').value;
      const primer_nombre = document.getElementById('primer_nombre').value;
      const segundo_nombre = document.getElementById('segundo_nombre').value;
      const primer_apellido = document.getElementById('primer_apellido').value;
      const segundo_apellido = document.getElementById('segundo_apellido').value;
      const email = document.getElementById('email').value;
      const telefono = document.getElementById('telefono').value;
      const fecha_nacimiento = document.getElementById('fecha_nacimiento').value;

      console.log('Form data:', { id_tipo_cuenta, tipo_documento, numero_documento, primer_nombre, fecha_nacimiento });

      fetch('/cuentas', {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: `id_tipo_cuenta=${encodeURIComponent(id_tipo_cuenta)}&tipo_documento=${encodeURIComponent(tipo_documento)}&numero_documento=${encodeURIComponent(numero_documento)}&primer_nombre=${encodeURIComponent(primer_nombre)}&segundo_nombre=${encodeURIComponent(segundo_nombre)}&primer_apellido=${encodeURIComponent(primer_apellido)}&segundo_apellido=${encodeURIComponent(segundo_apellido)}&email=${encodeURIComponent(email)}&telefono=${encodeURIComponent(telefono)}&fecha_nacimiento=${encodeURIComponent(fecha_nacimiento)}`
      }).then(async resp => {
        console.log('Response status:', resp.status);
        const alert = document.createElement('div');
        if (resp.ok) {
          console.log('Response OK');
          const data = await resp.json();
          alert.className = 'alert alert-success mt-2';
          alert.textContent = `${data.message}. Número: ${data.numero_cuenta}, ID: ${data.id_cuenta}`;
          form.appendChild(alert);
          setTimeout(() => alert.remove(), 4000);
          form.reset();
        } else {
          console.log('Response not OK');
          let txt = '';
          try { const data = await resp.json(); txt = data.message || JSON.stringify(data); } catch { txt = await resp.text(); }
          alert.className = 'alert alert-danger mt-2';
          alert.textContent = txt || 'No se pudo crear la cuenta';
          form.appendChild(alert);
          setTimeout(() => alert.remove(), 4000);
        }
      }).catch(error => {
        console.error('Fetch error:', error);
      });
    });
  } else {
    console.log('Form not found');
  }
});


