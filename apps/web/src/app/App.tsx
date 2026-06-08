import { Navigate, Route, Routes } from 'react-router-dom';
import { LoginPage } from '../features/auth/LoginPage';
import { DashboardPage } from '../features/dashboard/DashboardPage';
import { ClassesPage } from '../features/classes/ClassesPage';
import { ClassDetailPage } from '../features/classes/ClassDetailPage';
import { ActivityFormPage } from '../features/activities/ActivityFormPage';
import { Layout } from '../components/Layout';
import { getStoredTeacher } from '../lib/api';

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  if (!getStoredTeacher()) return <Navigate to="/login" replace />;
  return <>{children}</>;
}

export function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route
        path="/"
        element={
          <ProtectedRoute>
            <Layout />
          </ProtectedRoute>
        }
      >
        <Route index element={<DashboardPage />} />
        <Route path="classes" element={<ClassesPage />} />
        <Route path="classes/:id" element={<ClassDetailPage />} />
        <Route path="classes/:id/activities/new" element={<ActivityFormPage />} />
        <Route path="classes/:id/activities/:activityId/edit" element={<ActivityFormPage />} />
      </Route>
    </Routes>
  );
}
