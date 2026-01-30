# Pruebas Manuales — Sprint 1

## HU-CIT-001 — Agendar Citas

**Criterios de aceptacion:**
1. La recepcionista puede agendar una cita seleccionando paciente, fecha, hora y tipo de tratamiento
2. La cita queda registrada con estado inicial NUEVA
3. Se puede cambiar el estado de la cita segun el flujo: NUEVA > AGENDADA > CONFIRMADA > ATENDIDA/NO_ASISTIO/CANCELADA
4. Se puede reagendar una cita (cambia a estado REAGENDADA y se crea con nueva fecha/hora)

**Pasos:**

1. Iniciar sesion como **Administradora** (`admin@centrocaribel.com`)
2. Ir a **Citas** en el menu lateral
3. Hacer clic en **"Nueva Cita"**
4. Seleccionar un paciente, fecha, hora, tipo de tratamiento, turno (AM/PM) y opcionalmente observaciones
5. Guardar y verificar que la cita aparece en la lista con estado **NUEVA**
6. En la columna "Estado", abrir el selector y cambiar a **AGENDADA**. Verificar que se actualiza
7. Cambiar a **CONFIRMADA**, luego a **ATENDIDA**. Verificar que el flujo es respetado
8. Crear otra cita y probar el flujo: NUEVA > AGENDADA > **CANCELADA**
9. Crear otra cita, llevarla a AGENDADA, y hacer clic en **"Reagendar"**. Seleccionar nueva fecha/hora/turno. Verificar que el estado cambia a **REAGENDADA**
10. Probar con usuario **Interno** — NO deberia poder crear citas ni cambiar estados (botones no visibles)

---

## HU-PAC-001 — Registro de Pacientes

**Criterios de aceptacion:**
1. La recepcionista puede registrar un nuevo paciente con datos basicos
2. El sistema genera automaticamente un codigo unico para el paciente
3. Se puede buscar pacientes por nombre o CI

**Pasos:**

1. Iniciar sesion como **Administradora**
2. Ir a **Pacientes** en el menu lateral
3. Hacer clic en **"Nuevo Paciente"**
4. Llenar: nombre completo, CI, fecha de nacimiento, celular, direccion (opcional)
5. Guardar y verificar que el paciente aparece en la lista con un **codigo autogenerado** (ej: `PAC-0001`)
6. En el campo de busqueda, escribir parte del **nombre** del paciente. Verificar que filtra correctamente
7. Limpiar busqueda y buscar por **CI**. Verificar que encuentra al paciente
8. Probar con usuario **Interno** — puede ver la lista de pacientes pero NO deberia ver el boton "Nuevo Paciente"

---

## HU-PAC-002 — Gestionar Consentimiento Informado

**Criterios de aceptacion:**
1. La recepcionista puede registrar el consentimiento informado del paciente
2. El sistema registra el consentimiento con fecha, hora y firma del paciente
3. Se registra si el paciente autoriza o no el uso de fotos
4. El consentimiento queda vinculado al expediente del paciente

**Pasos:**

1. Iniciar sesion como **Administradora**
2. Ir a **Pacientes** y seleccionar un paciente existente (o crear uno nuevo)
3. En la seccion "Consentimientos Informados", hacer clic en **"Nuevo"**
4. Llenar el contenido del consentimiento, **firmar en el canvas** con el mouse/dedo, marcar o desmarcar "Autoriza fotos"
5. Hacer clic en "Registrar Consentimiento"
6. Verificar que el consentimiento aparece en la lista con:
   - Fecha y hora de registro
   - Badge "Autoriza fotos" o "No autoriza fotos"
   - El contenido ingresado
   - **La firma digital** del paciente (imagen visible)
7. Probar con un usuario **Interno** — NO deberia poder crear consentimientos (el boton "Nuevo" no aparece)

---

## HU-PAC-003 — Historia Clinica Digital

**Criterios de aceptacion:**
1. Al registrar un paciente, se crea automaticamente su historia clinica con numero secuencial
2. La historia incluye: antecedentes personales, familiares, alergias y medicamentos actuales
3. Se pueden registrar notas de evolucion con tipo (Tratamiento, Evolucion, Nota)
4. Solo personal autorizado puede editar la historia (Administradora, Licenciada para antecedentes; + Medico para notas)

**Pasos:**

1. **Crear un paciente nuevo** (como Administradora) y verificar que al entrar a su detalle aparece la seccion "Historia Clinica" con un numero secuencial (ej: `HC-0005`) y estado "ACTIVA"
2. **Registrar antecedentes**: hacer clic en "Registrar" en la seccion Historia Clinica, llenar los 4 campos (antecedentes personales, familiares, alergias, medicamentos actuales) y guardar. Verificar que aparecen en la tarjeta
3. **Editar antecedentes**: hacer clic en "Editar", modificar un campo y guardar. Verificar que se actualizo correctamente
4. **Crear nota de evolucion**: hacer clic en "Nueva Nota", seleccionar tipo "Evolucion", escribir contenido y guardar. Verificar que aparece en la lista con el badge correcto
5. **Crear nota de tratamiento**: repetir con tipo "Tratamiento". Verificar badge azul
6. **Crear nota general**: repetir con tipo "Nota". Verificar badge gris
7. **Probar con usuario Medico**: iniciar sesion como Medico, entrar al paciente:
   - Deberia poder **ver** la historia y antecedentes (solo lectura, sin boton "Editar")
   - Deberia poder **crear notas** de evolucion (boton "Nueva Nota" visible)
8. **Probar con usuario Interno**: iniciar sesion como Interno:
   - Deberia poder **ver** la historia, antecedentes y notas (solo lectura)
   - **NO** deberia ver boton "Editar" antecedentes ni "Nueva Nota"

---

## HU-ADM-001 — Gestion de Usuarios y Roles

**Criterios de aceptacion:**
1. La administradora puede crear usuarios con roles predefinidos
2. Roles disponibles: Administradora, Licenciada, Interno, Medico
3. Cada rol tiene permisos especificos de acceso al sistema
4. Se puede activar/desactivar usuarios

**Pasos:**

1. **Iniciar sesion como Administradora** (`admin@centrocaribel.com`)
2. Ir a **Usuarios** en el menu lateral
3. **Crear un usuario nuevo**: clic en "Nuevo Usuario", llenar nombre, email, contrasena, y seleccionar rol "Interno". Guardar y verificar que aparece en la lista
4. **Verificar los 4 roles disponibles**: al crear otro usuario, el selector debe mostrar exactamente: Administradora, Licenciada, Interno, Medico
5. **Desactivar un usuario**: en la lista, usar el boton de desactivar/toggle en un usuario. Verificar que su estado cambia a inactivo
6. **Verificar que usuario desactivado no puede ingresar**: cerrar sesion e intentar login con el usuario desactivado. Debe rechazar el acceso
7. **Reactivar el usuario**: volver a iniciar sesion como Administradora, activar al usuario nuevamente. Verificar que ahora si puede hacer login
8. **Verificar permisos por rol**:
   - Como **Licenciada**: puede ver pacientes, crear citas, registrar consentimientos, editar antecedentes, crear notas
   - Como **Interno**: solo puede ver pacientes, citas, historia (lectura). NO puede crear nada
   - Como **Medico**: puede ver todo + crear notas de evolucion. NO puede crear pacientes ni citas
9. **Verificar que solo Administradora ve la seccion "Usuarios"**: los otros roles NO deben ver el enlace "Usuarios" en el menu
