import React from 'react';
import { IStructuredError } from '../../interfaces/structured-error.interface';
import { Fatality } from '../fatality/fatality.component';

interface Props {
  message?: string;
  children: React.ReactNode;
}

interface State {
  hasError: boolean;
}

/**
 * Catches any unhandled errors that occur in props.children so that we can show an error message to the user. Without
 * this React will remove its UI from the screen instead. React error boundaries are only supported by class
 * components.
 *
 * See: https://react.dev/reference/react/Component#catching-rendering-errors-with-an-error-boundary
 */
export class ErrorBoundary extends React.Component<Props, State> {
  constructor(props: any) {
    super(props);
    this.state = { hasError: false };
  }

  /**
   * Lets us update state in response to an error.
   */
  static getDerivedStateFromError(): State {
    return { hasError: true };
  }

  /**
   * Called after an error has been thrown by a descendant component. Used here purely to log the
   * error, since getDerivedStateFromError cannot access it - without this, errors caught by the
   * boundary were silently swallowed.
   */
  componentDidCatch(error: Error, errorInfo: React.ErrorInfo): void {
    console.error('ErrorBoundary caught an error:', error, errorInfo.componentStack);
  }

  render(): React.ReactNode {
    const { children, message } = this.props;

    if (this.state.hasError) {
      const error = {} as IStructuredError;
      const fallback = message || 'An unexpected error has occurred';

      return <Fatality error={error} message={fallback} />;
    }

    return children;
  }
}
