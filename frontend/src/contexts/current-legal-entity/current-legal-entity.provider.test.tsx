import React from 'react';
import { fireEvent, render, waitFor } from '@testing-library/react';
import { screen } from '@testing-library/dom';
import '@testing-library/jest-dom';
import { MemoryRouter, Route, Routes, useNavigate } from 'react-router-dom';
import { CurrentLegalEntityProvider } from './current-legal-entity.provider';
import { CurrentLegalEntityContext } from './current-legal-entity.context';
import { LegalEntitiesContext } from '../legal-entities/legal-entities.context';
import { SelectedLegalEntityContext } from '../selected-legal-entity/selected-legal-entity.context';
import { ILegalEntity } from '../../interfaces/legal-entity.interface';

function makeLegalEntity(name: string): ILegalEntity {
  return { name } as ILegalEntity;
}

function ShowCurrentLegalEntity(): JSX.Element {
  const { currentLegalEntity } = React.useContext(CurrentLegalEntityContext);
  return <div>Current: {currentLegalEntity?.name}</div>;
}

function NavigateButton({ to }: { to: string }): JSX.Element {
  const navigate = useNavigate();
  return <button onClick={() => navigate(to)}>navigate</button>;
}

describe('Current legal entity provider', () => {
  it('re-fetches the legal entity when the route legal_entity_name param changes', async () => {
    const getLegalEntityByName = jest.fn((name: string) => Promise.resolve(makeLegalEntity(name)));

    render(
      <MemoryRouter initialEntries={['/orgs/foo']}>
        <LegalEntitiesContext.Provider value={{ legalEntities: [], getLegalEntityById: jest.fn(), getLegalEntityByName }}>
          <SelectedLegalEntityContext.Provider value={{ selectedLegalEntity: makeLegalEntity('foo'), selectLegalEntity: jest.fn() }}>
            <Routes>
              <Route
                path="/orgs/:legal_entity_name"
                element={
                  <>
                    <NavigateButton to="/orgs/bar" />
                    <CurrentLegalEntityProvider>
                      <ShowCurrentLegalEntity />
                    </CurrentLegalEntityProvider>
                  </>
                }
              />
            </Routes>
          </SelectedLegalEntityContext.Provider>
        </LegalEntitiesContext.Provider>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText('Current: foo')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('navigate'));

    await waitFor(() => {
      expect(screen.getByText('Current: bar')).toBeInTheDocument();
    });

    expect(getLegalEntityByName).toHaveBeenCalledWith('foo');
    expect(getLegalEntityByName).toHaveBeenCalledWith('bar');
  });
});
