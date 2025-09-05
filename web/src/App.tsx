import { Routes, Route } from 'react-router-dom';
import { Layout } from '@/components/layout/Layout';
import { Dashboard } from '@/pages/Dashboard';
import { Connections } from '@/pages/Connections';
import { Analysis } from '@/pages/Analysis';
import { CollectionView } from '@/pages/CollectionView';

const Security = () => <div className="p-6"><h1 className="text-3xl font-bold">Security Analysis</h1><p className="text-onyx mt-2">Coming soon...</p></div>;
const Optimization = () => <div className="p-6"><h1 className="text-3xl font-bold">Query Optimization</h1><p className="text-onyx mt-2">Coming soon...</p></div>;
const Monitoring = () => <div className="p-6"><h1 className="text-3xl font-bold">Monitoring & Alerts</h1><p className="text-onyx mt-2">Coming soon...</p></div>;
const Lineage = () => <div className="p-6"><h1 className="text-3xl font-bold">Data Lineage</h1><p className="text-onyx mt-2">Coming soon...</p></div>;
const Settings = () => <div className="p-6"><h1 className="text-3xl font-bold">Settings</h1><p className="text-onyx mt-2">Coming soon...</p></div>;

function App() {
  return (
    <Routes>
      <Route path="/" element={<Layout />}>
        <Route index element={<Dashboard />} />
        <Route path="connections" element={<Connections />} />
        <Route path="analysis" element={<Analysis />} />
        <Route path="collection/:connectionId/:collectionName" element={<CollectionView />} />
        <Route path="security" element={<Security />} />
        <Route path="optimization" element={<Optimization />} />
        <Route path="monitoring" element={<Monitoring />} />
        <Route path="lineage" element={<Lineage />} />
        <Route path="settings" element={<Settings />} />
      </Route>
    </Routes>
  );
}

export default App;
