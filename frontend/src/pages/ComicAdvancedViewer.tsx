import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Stage, Layer, Image as KonvaImage, Rect, Text } from 'react-konva';
import useImage from 'use-image';
import { generationApi } from '../api/generations';
import { ArrowLeft, Save } from 'lucide-react';

const URLImage = ({ src, x, y, width, height, isSplash }: any) => {
  const [img] = useImage(src);
  return (
    <KonvaImage
      image={img}
      x={x}
      y={y}
      width={isSplash ? width * 2 + 20 : width}
      height={isSplash ? height * 1.5 : height}
      draggable
      stroke="#ccc"
      strokeWidth={2}
    />
  );
};

export const ComicAdvancedViewer = () => {
  const { genId } = useParams<{ genId: string }>();
  const navigate = useNavigate();
  const [gen, setGen] = useState<any>(null);
  const [rawUrls, setRawUrls] = useState<string[]>([]);
  const [showText, setShowText] = useState(true);
  const stageRef = React.useRef<any>(null);

  const handleExportAyk = () => {
    if (!stageRef.current) return;
    const jsonStr = stageRef.current.toJSON();
    const blob = new Blob([jsonStr], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `comic_${genId || 'project'}.ayk`;
    a.click();
    URL.revokeObjectURL(url);
  };

  const handleImportAyk = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    const reader = new FileReader();
    reader.onload = (event) => {
      const jsonStr = event.target?.result as string;
      if (jsonStr) {
        alert("Project Tree Hydrated:\n" + jsonStr.substring(0, 50) + "...\n(In full implementation, this sets React state to mirror layers)");
      }
    };
    reader.readAsText(file);
  };
  
  useEffect(() => {
    if (genId) {
      generationApi.get(genId).then(res => {
        setGen(res.data);
        if (res.data.raw_image_urls) {
          try { setRawUrls(JSON.parse(res.data.raw_image_urls)); } catch(e) {}
        }
      });
    }
  }, [genId]);

  if (!gen) return <div style={{ padding: 40, textAlign: 'center' }}>Loading Canvas...</div>;

  return (
    <div style={{ display: 'flex', flexDirection: 'column', height: '100vh', background: '#ececec' }}>
      <header style={{ padding: '15px 20px', background: '#fff', borderBottom: '1px solid #ddd', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
         <button onClick={() => navigate('/dashboard')} style={{ display: 'flex', alignItems: 'center', gap: 5, padding: '8px 12px', background: 'transparent', border: '1px solid #ccc', borderRadius: 4, cursor: 'pointer' }}>
          <ArrowLeft size={16} /> Dashboard
        </button>
        <h3 style={{ margin: 0 }}>Advanced Composer (Konva Canvas)</h3>
        <div style={{ display: 'flex', gap: 15, alignItems: 'center' }}>
          <label style={{ display: 'flex', alignItems: 'center', gap: 5, fontSize: '0.9em', color: '#555', cursor: 'pointer' }}>
            <input type="checkbox" checked={showText} onChange={e => setShowText(e.target.checked)} /> Show Text Layer
          </label>
          <div style={{ borderLeft: '1px solid #ccc', height: 20 }}></div>
          <button onClick={handleExportAyk} style={{ display: 'flex', alignItems: 'center', gap: 5, padding: '8px 12px', background: '#333', color: '#fff', border: 'none', borderRadius: 4, cursor: 'pointer' }}>
            Save .ayk
          </button>
          <label style={{ display: 'flex', alignItems: 'center', gap: 5, padding: '8px 12px', background: '#ccc', color: '#000', border: 'none', borderRadius: 4, cursor: 'pointer' }}>
            Load .ayk
            <input type="file" accept=".ayk,.json" style={{ display: 'none' }} onChange={handleImportAyk} />
          </label>
          <button style={{ display: 'flex', alignItems: 'center', gap: 5, padding: '8px 16px', background: '#10b981', color: '#fff', border: 'none', borderRadius: 4, cursor: 'pointer', fontWeight: 'bold' }}>
            <Save size={16} /> Export Render
          </button>
        </div>
      </header>

      <div style={{ flex: 1, display: 'flex', justifyContent: 'center', alignItems: 'center', overflow: 'hidden' }}>
        <div style={{ background: '#fff', boxShadow: '0 8px 30px rgba(0,0,0,0.1)' }}>
          <Stage width={800} height={900} ref={stageRef}>
            <Layer>
              <Rect width={800} height={900} fill="#ffffff" />
              
              {rawUrls.map((url, idx) => {
                const isSplash = idx === 1; // Example: Mock standard dynamic splash on panel #2
                let row = Math.floor(idx / 2);
                if (idx > 1) row = Math.floor((idx + 1) / 2); // Shift logic purely for splash demo visualization 
                
                const col = idx % 2;
                const size = 350;
                const padding = 30;
                
                return (
                  <URLImage 
                    key={idx} 
                    src={url} 
                    x={padding + col * (size + padding)} 
                    y={padding + row * (size + padding)} 
                    width={size} 
                    height={size} 
                    isSplash={isSplash}
                  />
                );
              })}
              
              <Text 
                text="Drag Me! Overlapping Layer" 
                x={400} 
                y={500} 
                fontSize={32} 
                fontFamily="Impact, sans-serif"
                fill="#fff"
                draggable
                stroke="#000"
                strokeWidth={3}
                shadowColor="rgba(0,0,0,0.5)"
                shadowBlur={5}
                shadowOffset={{ x: 2, y: 2 }}
                shadowOpacity={0.8}
                visible={showText}
              />
            </Layer>
          </Stage>
        </div>
      </div>
    </div>
  );
};
