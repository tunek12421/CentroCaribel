import { useRef, useState } from 'react';
import { Button } from '../../components/ui/Button';
import type { CreateConsentimientoRequest } from '../../types';

interface ConsentimientoFormProps {
  onSubmit: (data: CreateConsentimientoRequest) => void;
  loading: boolean;
  error: unknown;
}

export function ConsentimientoForm({ onSubmit, loading, error }: ConsentimientoFormProps) {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const [isDrawing, setIsDrawing] = useState(false);
  const [contenido, setContenido] = useState('');
  const [autorizaFotos, setAutorizaFotos] = useState(false);
  const [hasSigned, setHasSigned] = useState(false);

  const apiError = (error as any)?.response?.data?.error;

  const getPos = (e: React.MouseEvent | React.TouchEvent) => {
    const canvas = canvasRef.current!;
    const rect = canvas.getBoundingClientRect();
    if ('touches' in e) {
      return { x: e.touches[0].clientX - rect.left, y: e.touches[0].clientY - rect.top };
    }
    return { x: e.clientX - rect.left, y: e.clientY - rect.top };
  };

  const startDraw = (e: React.MouseEvent | React.TouchEvent) => {
    e.preventDefault();
    const ctx = canvasRef.current?.getContext('2d');
    if (!ctx) return;
    const pos = getPos(e);
    ctx.beginPath();
    ctx.moveTo(pos.x, pos.y);
    setIsDrawing(true);
    setHasSigned(true);
  };

  const draw = (e: React.MouseEvent | React.TouchEvent) => {
    e.preventDefault();
    if (!isDrawing) return;
    const ctx = canvasRef.current?.getContext('2d');
    if (!ctx) return;
    const pos = getPos(e);
    ctx.lineWidth = 2;
    ctx.lineCap = 'round';
    ctx.strokeStyle = '#1a1a1a';
    ctx.lineTo(pos.x, pos.y);
    ctx.stroke();
  };

  const stopDraw = () => setIsDrawing(false);

  const clearCanvas = () => {
    const ctx = canvasRef.current?.getContext('2d');
    if (ctx && canvasRef.current) {
      ctx.clearRect(0, 0, canvasRef.current.width, canvasRef.current.height);
    }
    setHasSigned(false);
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    let firma: string | undefined;
    if (hasSigned && canvasRef.current) {
      const dataUrl = canvasRef.current.toDataURL('image/png');
      firma = dataUrl.split(',')[1]; // base64 sin header
    }
    onSubmit({ firma_digital: firma, autoriza_fotos: autorizaFotos, contenido });
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {apiError && (
        <div className="p-3 bg-red-50 border border-red-200 rounded-lg text-sm text-danger">
          {apiError.detail || apiError.message}
        </div>
      )}

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">Contenido del consentimiento *</label>
        <textarea
          value={contenido}
          onChange={(e) => setContenido(e.target.value)}
          rows={4}
          className="w-full"
          placeholder="Yo, autorizo el tratamiento de fisioterapia estética..."
          required
        />
      </div>

      <div className="flex items-center gap-2">
        <input
          type="checkbox"
          id="autorizaFotos"
          checked={autorizaFotos}
          onChange={(e) => setAutorizaFotos(e.target.checked)}
          className="h-4 w-4 rounded border-border text-primary focus:ring-primary"
        />
        <label htmlFor="autorizaFotos" className="text-sm text-gray-700">
          Autoriza toma de fotografías
        </label>
      </div>

      <div>
        <div className="flex items-center justify-between mb-1">
          <label className="text-sm font-medium text-gray-700">Firma digital</label>
          <button type="button" onClick={clearCanvas} className="text-xs text-primary hover:underline">
            Limpiar
          </button>
        </div>
        <canvas
          ref={canvasRef}
          width={400}
          height={150}
          className="w-full border border-border rounded-lg bg-white cursor-crosshair touch-none"
          onMouseDown={startDraw}
          onMouseMove={draw}
          onMouseUp={stopDraw}
          onMouseLeave={stopDraw}
          onTouchStart={startDraw}
          onTouchMove={draw}
          onTouchEnd={stopDraw}
        />
      </div>

      <Button type="submit" loading={loading} className="w-full" disabled={!contenido.trim()}>
        Registrar Consentimiento
      </Button>
    </form>
  );
}
