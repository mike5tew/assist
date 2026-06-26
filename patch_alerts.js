const fs = require('fs');
const file = '/Users/michaelstewart/Coding/SampleTrack/ReactNativeMobileFrontEnd/src/ScanScreen.tsx';
let data = fs.readFileSync(file, 'utf8');

// Replace Alert in handleLocationBarcodeScan
data = data.replace(
  /Alert\.alert\('Unknown Location', `No location found for barcode: \$\{data\}`\);/g,
  "showToast(`Unknown location: ${data}`);"
);

// In handleBoxScan create box catch
data = data.replace(
  /\} catch \(e: any\) \{\s*Alert\.alert\('Error', e\.message\);\s*\}/g,
  `} catch (e: any) {\n        showToast(e.message || 'Error configuring box');\n      }`
);

// Add lock to handleBoxBarcodeScan
data = data.replace(
  /const handleBoxBarcodeScan = \(\{ data \}: \{ data: string \}\) => \{\s*if \(!scanningBoxBarcode\) return;\s*setScanningBoxBarcode\(false\);/g,
  `const handleBoxBarcodeScan = async ({ data }: { data: string }) => {\n    if (!scanningBoxBarcode || isScanningRef.current) return;\n    isScanningRef.current = true;\n    setScanningBoxBarcode(false);\n    const unlock = () => setTimeout(() => { isScanningRef.current = false; }, 800);`
);

// We must also call unlock at the end of handleBoxBarcodeScan
data = data.replace(
  /setupBox\(data\);\s*\};/g,
  `// Note: setupBox handles its own await now, but wait, setupBox is async
    await setupBox(data);
    unlock();
  };`
);

// We need to modify setupBox to be async and catch errors to unlock, OR just rely on showToast
// Actually, earlier we saw setupBox was just an async function? Let's look.

fs.writeFileSync(file, data);
