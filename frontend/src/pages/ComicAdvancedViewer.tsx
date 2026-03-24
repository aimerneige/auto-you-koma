import { useState, useEffect } from 'react';
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
        <button style={{ display: 'flex', alignItems: 'center', gap: 5, padding: '8px 16px', background: '#10b981', color: '#fff', border: 'none', borderRadius: 4, cursor: 'pointer', fontWeight: 'bold' }}>
          <Save size={16} /> Export Render
        </button>
      </header>

      <div style={{ flex: 1, display: 'flex', justifyContent: 'center', alignItems: 'center', overflow: 'hidden' }}>
        <div style={{ background: '#fff', boxShadow: '0 8px 30px rgba(0,0,0,0.1)' }}>
          <Stage width={800} height={900}>
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
              />
            </Layer>
          </Stage>
        </div>
      </div>
    </div>
  );
};
