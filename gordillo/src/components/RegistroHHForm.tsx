import React, { useState } from 'react';
import { Plus, Calendar, User, Briefcase, Clock } from 'lucide-react';

interface TrabajadorHH {
  trabajador: string;
  horasNormales: number;
  horasExtras: number;
}

interface RegistroHHData {
  fechaEjecucion: string;
  orderDeTrabajo: string;
  actividad: string;
  trabajadores: TrabajadorHH[];
}

interface RegistroHHFormProps {
  onSubmit: (registro: RegistroHHData) => void;
}

const RegistroHHForm: React.FC<RegistroHHFormProps> = ({ onSubmit }) => {
  const [formData, setFormData] = useState<RegistroHHData>({
    fechaEjecucion: '',
    orderDeTrabajo: '',
    actividad: '',
    trabajadores: [{ trabajador: '', horasNormales: 0, horasExtras: 0 }]
  });

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleTrabajadorChange = (index: number, field: keyof TrabajadorHH, value: string | number) => {
    setFormData(prev => ({
      ...prev,
      trabajadores: prev.trabajadores.map((trabajador, i) => 
        i === index ? { ...trabajador, [field]: value } : trabajador
      )
    }));
  };

  const addTrabajador = () => {
    setFormData(prev => ({
      ...prev,
      trabajadores: [...prev.trabajadores, { trabajador: '', horasNormales: 0, horasExtras: 0 }]
    }));
  };

  const removeTrabajador = (index: number) => {
    if (formData.trabajadores.length > 1) {
      setFormData(prev => ({
        ...prev,
        trabajadores: prev.trabajadores.filter((_, i) => i !== index)
      }));
    }
  };

  const addExtraHour = (index: number) => {
    setFormData(prev => ({
      ...prev,
      trabajadores: prev.trabajadores.map((trabajador, i) => 
        i === index ? { ...trabajador, horasExtras: trabajador.horasExtras + 1 } : trabajador
      )
    }));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const hasValidTrabajadores = formData.trabajadores.some(t => t.trabajador && (t.horasNormales > 0 || t.horasExtras > 0));
    
    if (formData.fechaEjecucion && formData.actividad && hasValidTrabajadores) {
      onSubmit(formData);
      setFormData({
        fechaEjecucion: '',
        orderDeTrabajo: '',
        actividad: '',
        trabajadores: [{ trabajador: '', horasNormales: 0, horasExtras: 0 }]
      });
    }
  };

  return (
    <div className="bg-white rounded-xl shadow-xl p-8 max-w-4xl mx-auto border border-gray-100">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-800 mb-2">Registro HH</h1>
        <p className="text-gray-600 text-lg">Registro de Horas Trabajador</p>
      </div>

      <form onSubmit={handleSubmit} className="space-y-6">
        {/* Información general */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div className="space-y-2">
            <label htmlFor="fechaEjecucion" className="block text-sm font-semibold text-gray-700">
              Fecha
            </label>
            <div className="relative">
              <Calendar className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
              <input
                type="date"
                id="fechaEjecucion"
                name="fechaEjecucion"
                value={formData.fechaEjecucion}
                onChange={handleInputChange}
                className="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200 text-gray-700"
                required
              />
            </div>
          </div>

          <div className="space-y-2">
            <label htmlFor="orderDeTrabajo" className="block text-sm font-semibold text-gray-700">
              Orden de Trabajo
            </label>
            <div className="relative">
              <Briefcase className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
              <select
                id="orderDeTrabajo"
                name="orderDeTrabajo"
                value={formData.orderDeTrabajo}
                onChange={handleInputChange}
                className="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200 text-gray-700 appearance-none bg-white"
              >
                <option value="">N/A</option>
                <option value="OT-001">OT-001</option>
                <option value="OT-002">OT-002</option>
                <option value="OT-003">OT-003</option>
                <option value="OT-004">OT-004</option>
              </select>
            </div>
          </div>

          <div className="space-y-2">
            <label htmlFor="actividad" className="block text-sm font-semibold text-gray-700">
              Actividad
            </label>
            <div className="relative">
              <Clock className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
              <select
                id="actividad"
                name="actividad"
                value={formData.actividad}
                onChange={handleInputChange}
                className="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200 text-gray-700 appearance-none bg-white"
                required
              >
                <option value="">Seleccionar...</option>
                <option value="trabajo de mantenimiento">Trabajo de Mantenimiento</option>
                <option value="produccion">Producción</option>
                <option value="calidad">Control de Calidad</option>
                <option value="logistica">Logística</option>
                <option value="administrativo">Administrativo</option>
              </select>
            </div>
          </div>
        </div>

        {/* Trabajadores */}
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-semibold text-gray-800">Trabajadores</h3>
            <button
              type="button"
              onClick={addTrabajador}
              className="bg-green-600 text-white px-4 py-2 rounded-lg hover:bg-green-700 transition-colors duration-200 flex items-center space-x-2"
            >
              <Plus className="w-4 h-4" />
              <span>Agregar Trabajador</span>
            </button>
          </div>

          {formData.trabajadores.map((trabajador, index) => (
            <div key={index} className="bg-gray-50 p-4 rounded-lg border border-gray-200">
              <div className="grid grid-cols-1 md:grid-cols-5 gap-4 items-end">
                <div className="space-y-2">
                  <label className="block text-sm font-semibold text-gray-700">
                    Trabajador
                  </label>
                  <div className="relative">
                    <User className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
                    <select
                      value={trabajador.trabajador}
                      onChange={(e) => handleTrabajadorChange(index, 'trabajador', e.target.value)}
                      className="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200 text-gray-700 appearance-none bg-white"
                      required
                    >
                      <option value="">Seleccionar...</option>
                      <option value="alexander gutierrez">Alexander Gutierrez</option>
                      <option value="maria rodriguez">María Rodríguez</option>
                      <option value="carlos lopez">Carlos López</option>
                      <option value="ana martinez">Ana Martínez</option>
                      <option value="juan perez">Juan Pérez</option>
                      <option value="lucia garcia">Lucía García</option>
                    </select>
                  </div>
                </div>

                <div className="space-y-2">
                  <label className="block text-sm font-semibold text-gray-700">
                    H. Normales
                  </label>
                  <input
                    type="number"
                    value={trabajador.horasNormales}
                    onChange={(e) => handleTrabajadorChange(index, 'horasNormales', parseFloat(e.target.value) || 0)}
                    className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200 text-gray-700"
                    placeholder="0"
                    min="0"
                    step="0.5"
                  />
                </div>

                <div className="space-y-2">
                  <label className="block text-sm font-semibold text-gray-700">
                    H. Extras
                  </label>
                  <div className="flex">
                    <input
                      type="number"
                      value={trabajador.horasExtras}
                      onChange={(e) => handleTrabajadorChange(index, 'horasExtras', parseFloat(e.target.value) || 0)}
                      className="flex-1 px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200 text-gray-700"
                      placeholder="0"
                      min="0"
                      step="0.5"
                    />
                    <button
                      type="button"
                      onClick={() => addExtraHour(index)}
                      className="ml-2 bg-blue-600 text-white px-3 py-3 rounded-lg hover:bg-blue-700 transition-colors duration-200 flex items-center justify-center"
                    >
                      <Plus className="w-5 h-5" />
                    </button>
                  </div>
                </div>

                <div className="space-y-2">
                  <label className="block text-sm font-semibold text-gray-700">
                    Total Horas
                  </label>
                  <div className="px-4 py-3 bg-gray-100 border border-gray-300 rounded-lg text-center font-semibold text-gray-700">
                    {trabajador.horasNormales + trabajador.horasExtras}h
                  </div>
                </div>

                <div className="space-y-2">
                  <label className="block text-sm font-semibold text-gray-700">
                    Acciones
                  </label>
                  {formData.trabajadores.length > 1 ? (
                    <button
                      type="button"
                      onClick={() => removeTrabajador(index)}
                      className="w-full bg-red-600 text-white px-4 py-3 rounded-lg hover:bg-red-700 transition-colors duration-200"
                    >
                      Eliminar
                    </button>
                  ) : (
                    <div className="w-full px-4 py-3 bg-gray-200 border border-gray-300 rounded-lg text-center text-gray-500 text-sm">
                      Mínimo 1
                    </div>
                  )}
                </div>
              </div>
            </div>
          ))}
        </div>

        <div className="flex justify-center pt-6">
          <button
            type="submit"
            className="bg-blue-600 text-white px-12 py-4 rounded-lg font-semibold hover:bg-blue-700 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 shadow-lg hover:shadow-xl transform hover:-translate-y-0.5"
          >
            CREAR REGISTRO
          </button>
        </div>
      </form>

      <div className="mt-8 text-center">
        <button className="text-blue-600 hover:text-blue-800 font-medium transition-colors duration-200">
          Volver a la Lista
        </button>
      </div>
    </div>
  );
};

export default RegistroHHForm;
