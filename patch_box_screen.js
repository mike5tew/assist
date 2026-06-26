const fs = require('fs');
const file = '/Users/michaelstewart/Coding/SampleTrack/ReactNativeMobileFrontEnd/src/BoxScreen.tsx';
let data = fs.readFileSync(file, 'utf8');

// Update imports if needed (ensure useRef etc exists)
if (!data.includes('useRef')) {
  data = data.replace(/useState/g, 'useState, useRef');
}

// Add showToast and isScanningRef
const toastState = `
  const isScanningRef = useRef(false);
  const [toastMessage, setToastMessage] = useState<string | null>(null);

  const showToast = useCallback((msg: string) => {
    setToastMessage(msg);
    setTimeout(() => {
      setToastMessage(prev => prev === msg ? null : prev);
    }, 2500);
  }, []);
`;

data = data.replace(
  /const \[scanned, setScanned\] = useState\(false\);/g,
  `const [scanned, setScanned] = useState(false);\n${toastState}`
);

// We also need to mount the toast in the UI
if (!data.includes('toastContainer')) {
  data = data.replace(
    /<View style=\{styles\.container\}>/g,
    `<View style={styles.container}>
      {toastMessage && (
        <View style={styles.toastContainer}>
          <Text style={styles.toastText}>{toastMessage}</Text>
        </View>
      )}`
  );
  
  // Add Toast styles
  data = data.replace(
    /container: \{/g,
    `toastContainer: {
    position: 'absolute',
    top: 60,
    left: 20,
    right: 20,
    backgroundColor: '#323232',
    padding: 16,
    borderRadius: 8,
    zIndex: 9999,
    flexDirection: 'row',
    alignItems: 'center',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.25,
    shadowRadius: 3.84,
    elevation: 5,
  },
  toastText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
    flex: 1,
  },
  container: {`
  );
}

// Refactor handleBarCodeScanned
data = data.replace(
  /const handleBarCodeScanned = async \(\{ data \}: \{ data: string \}\) => \{[\s\S]*?setScanned\(false\);\s*\};/g,
  `const handleBarCodeScanned = useCallback(async ({ data }: { data: string }) => {
    if (isScanningRef.current || scanned) return;
    
    isScanningRef.current = true;
    setScanned(true);

    const unlock = (delay = 800) => {
      setTimeout(() => {
        isScanningRef.current = false;
        setScanned(false);
      }, delay);
    };

    if (scanMode === 'box') {
      // Set active box
      if (!token) { unlock(0); return; }
      try {
        const box = await api.getBox(token, data);
        if (box.status !== 'open') {
          Alert.alert(
            'Reuse Box?',
            \`Box "\${box.label}" is sealed. Start a new collection in this box?\`,
            [
              { text: 'Cancel', style: 'cancel', onPress: () => unlock(0) },
              {
                text: 'Reuse',
                style: 'destructive',
                onPress: async () => {
                  try {
                    const reopened = await api.reopenBox(token, data);
                    setActiveBox(reopened);
                    loadBoxes();
                    unlock(0);
                  } catch (e: any) {
                    showToast(e.message);
                    unlock(2000);
                  }
                },
              },
            ],
          );
          return; // Wait for alert
        } else {
          setActiveBox(box);
        }
      } catch {
        // Not found -> ask to create or switch mode
        Alert.alert(
          'Box Not Found',
          \`Create new box "\${data}"?\`,
          [
            { text: 'Cancel', style: 'cancel', onPress: () => unlock(0) },
            {
              text: 'Create',
              onPress: async () => {
                try {
                  const box = await api.createBox(token, data);
                  setActiveBox(box);
                  loadBoxes();
                  unlock(0);
                } catch (e: any) {
                  showToast(e.message);
                  unlock(2000);
                }
              },
            },
          ],
        );
        return; // wait for alert
      }
    } else if (scanMode === 'sample' && activeBox) {
      if (!token) { unlock(0); return; }
      try {
        const updated = await api.addSampleToBox(token, activeBox.barcode, data);
        setActiveBox(updated);
        loadBoxes();
        showToast(\`Sample \${data} added to box\`);
      } catch (e: any) {
        showToast(e.message || 'Error adding sample');
        unlock(2000);
        return;
      }
    }

    unlock(800);
  }, [scanned, scanMode, activeBox, token, showToast]);`
);

fs.writeFileSync(file, data);
