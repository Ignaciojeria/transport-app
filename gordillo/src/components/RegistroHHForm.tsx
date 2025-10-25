import React, { useState } from 'react';
import { Plus, Calendar, User, Briefcase, Clock } from 'lucide-react';

interface RegistroHHData {
  fechaEjecucion: string;
  trabajador: string;
  orderDeTrabajo: string;
  horasNormales: number;
  horasExtras: number;
  actividad: string;
}

interface RegistroHHFormProps {
  onSubmit: (registro: RegistroHHData) => void;
}

const RegistroHHForm: React.FC<RegistroHHFormProps> = ({ onSubmit }) => {
  const [formData, setFormData] = useState<RegistroHHData>({
    fechaEjecucion: '',
    trabajador: '',
    orderDeTrabajo: '',
    horasNormales: 0,
    horasExtras: 0,
    actividad: ''
  });

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: name === 'horasNormales' || name === 'horasExtras' ? parseFloat(value) || 0 : value
    }));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (formData.fechaEjecucion && formData.trabajador && formData.actividad) {
      onSubmit(formData);
      setFormData({
        fechaEjecucion: '',
        trabajador: '',
        orderDeTrabajo: '',
        horasNormales: 0,
        horasExtras: 0,
        actividad: ''
      });
    }
  };

  const addExtraHour = () => {
    setFormData(prev => ({
      ...prev,
      horasExtras: prev.horasExtras + 1
    }));
  };

  return (
    <div className="bg-white rounded-xl shadow-xl p-8 max-w-4xl mx-auto border border-gray-100">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-800 mb-2">Registro HH</h1>
        <p className="text-gray-600 text-lg">Registro de Horas Trabajador</p>
      </div>

      <form onSubmit={handleSubmit} className="space-y-6">
        {/* Primera fila */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
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
            <label htmlFor="trabajador" className="block text-sm font-semibold text-gray-700">
              Trabajador
            </label>
            <div className="relative">
              <User className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
              <select
                id="trabajador"
                name="trabajador"
                value={formData.trabajador}
                onChange={handleInputChange}
                className="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200 text-gray-700 appearance-none bg-white"
                required
              >
                <option value="">Seleccionar...</option>
                <option value="alexander gutierrez">Alexander Gutierrez</option>
                <option value="maria rodriguez">María Rodríguez</option>
                <option value="carlos lopez">Carlos López</option>
                <option value="ana martinez">Ana Martínez</option>
              </select>
            </div>
          </div>
        </div>

        {/* Segunda fila */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-4">
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

          <div className="space-y-2">
            <label htmlFor="horasNormales" className="block text-sm font-semibold text-gray-700">
              H. Normales
            </label>
            <input
              type="number"
              id="horasNormales"
              name="horasNormales"
              value={formData.horasNormales}
              onChange={handleInputChange}
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200 text-gray-700"
              placeholder="0"
              min="0"
              step="0.5"
            />
          </div>

          <div className="space-y-2">
            <label htmlFor="horasExtras" className="block text-sm font-semibold text-gray-700">
              H. Extras
            </label>
            <div className="flex">
              <input
                type="number"
                id="horasExtras"
                name="horasExtras"
                value={formData.horasExtras}
                onChange={handleInputChange}
                className="flex-1 px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200 text-gray-700"
                placeholder="0"
                min="0"
                step="0.5"
              />
              <button
                type="button"
                onClick={addExtraHour}
                className="ml-2 bg-blue-600 text-white px-3 py-3 rounded-lg hover:bg-blue-700 transition-colors duration-200 flex items-center justify-center"
              >
                <Plus className="w-5 h-5" />
              </button>
            </div>
          </div>
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
